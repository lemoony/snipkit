package execution

import (
	"bytes"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	"emperror.dev/errors"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/creack/pty"
	"github.com/phuslu/log"
	"golang.org/x/term"
)

// Buffer sizes for reading from PTY and stdin.
const (
	ptyReadBufferSize   = 4096
	stdinReadBufferSize = 256
)

// Timing constants.
const (
	spinnerTickInterval   = 100 * time.Millisecond
	terminalSettleTimeout = 50 * time.Millisecond
)

// ctrlC is the ASCII code for Ctrl+C.
const ctrlC = 3

// controlSeqRegex matches terminal control sequences that shouldn't be stored in history.
var controlSeqRegex = regexp.MustCompile(`\x1b\][^\x07]*(?:\x07|\x1b\\)|\x1b\[[0-9;]*R`)

// csiRegex matches CSI sequences (cursor movement, colors, etc.).
var csiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[A-Za-z]`)

// CapturedOutput holds the result of script execution.
type CapturedOutput struct {
	Stdout   string
	ExitCode int
	Duration time.Duration
}

// helpLine renders a styled help line.
func helpLine(running bool, fromAssistant bool) string {
	gray := "\x1b[38;5;240m"
	reset := "\x1b[0m"

	if running {
		return gray + "Ctrl+C: abort • Script running..." + reset
	}
	if fromAssistant {
		return gray + "Enter: back to assistant • Ctrl+C: quit" + reset
	}
	return gray + "Enter/Ctrl+C: quit" + reset
}

// RunWithViewer executes the command with real-time output.
// When fromAssistant is true, shows UI elements (header, spinner, help line, waits for Enter).
// When fromAssistant is false, runs the command directly without any UI chrome.
func RunWithViewer(cmd *exec.Cmd, fromAssistant bool) *CapturedOutput {
	// Get terminal size
	cols, rows := 80, 24
	if w, h, sizeErr := term.GetSize(int(os.Stdout.Fd())); sizeErr == nil {
		cols, rows = w, h
	}

	// Start command in PTY
	ptmx, err := pty.Start(cmd)
	if err != nil {
		panic(errors.Wrapf(errors.WithStack(err), "failed to start pty"))
	}
	// Note: ptmx.Close() is called explicitly after command finishes, not deferred

	// Buffer to capture output
	var outputBuf bytes.Buffer

	// Channel to signal process completion
	done := make(chan struct{})

	// Direct execution: simple PTY without UI elements
	if !fromAssistant {
		return runDirect(cmd, ptmx, &outputBuf, done, cols, rows)
	}

	// Assistant execution: full UI with spinner, help line, and Enter wait
	return runWithAssistantUI(cmd, ptmx, &outputBuf, done, cols, rows)
}

// runDirect executes the command with PTY but without any UI elements.
// Used for direct command execution (not from assistant).
//
//nolint:gocognit,funlen // Complex function managing PTY and concurrent I/O
func runDirect(cmd *exec.Cmd, ptmx *os.File, outputBuf *bytes.Buffer, done chan struct{}, cols, rows int) *CapturedOutput {
	// Set PTY size
	_ = pty.Setsize(ptmx, &pty.Winsize{Rows: uint16(rows), Cols: uint16(cols)})

	// Handle window resize
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGWINCH)
	go func() {
		for range sigCh {
			if w, h, sizeErr := term.GetSize(int(os.Stdout.Fd())); sizeErr == nil {
				_ = pty.Setsize(ptmx, &pty.Winsize{Rows: uint16(h), Cols: uint16(w)})
			}
		}
	}()
	defer signal.Stop(sigCh)

	// Set terminal to raw mode for proper input handling
	oldState, rawErr := term.MakeRaw(int(os.Stdin.Fd()))
	rawModeEnabled := rawErr == nil

	// Forward stdin to PTY
	go func() {
		buf := make([]byte, stdinReadBufferSize)
		for {
			n, readErr := os.Stdin.Read(buf)
			if readErr != nil {
				break
			}
			if n > 0 {
				if _, writeErr := ptmx.Write(buf[:n]); writeErr != nil {
					return
				}
			}
		}
	}()

	// Copy output from PTY to stdout and buffer
	go func() {
		buf := make([]byte, ptyReadBufferSize)
		for {
			n, readErr := ptmx.Read(buf)
			if n > 0 {
				_, _ = os.Stdout.Write(buf[:n])
				outputBuf.Write(buf[:n])
			}
			if readErr != nil {
				break
			}
		}
		close(done)
	}()

	// Track execution time
	startTime := time.Now()

	// Wait for command to finish
	cmdErr := cmd.Wait()
	duration := time.Since(startTime)

	// Wait for output goroutine to finish
	<-done

	// Extract exit code
	exitCode := 0
	if cmdErr != nil {
		if exitErr, ok := cmdErr.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}

	// Close PTY
	_ = ptmx.Close()

	// Restore terminal from raw mode
	if rawModeEnabled {
		if restoreErr := term.Restore(int(os.Stdin.Fd()), oldState); restoreErr != nil {
			log.Warn().Err(restoreErr).Msg("Failed to restore terminal state")
		}
	}

	// Clean output for history
	rawOutput := outputBuf.String()
	outputToReturn := cleanOutput(rawOutput)
	if exitCode != 0 {
		outputToReturn = rawOutput
	}

	return &CapturedOutput{
		Stdout:   outputToReturn,
		ExitCode: exitCode,
		Duration: duration,
	}
}

// runWithAssistantUI executes the command with full UI elements (header, spinner, help line).
// Used for execution from the assistant.
//
//nolint:gocognit,gocyclo,funlen // Complex function managing PTY, terminal state, and concurrent I/O
func runWithAssistantUI(cmd *exec.Cmd, ptmx *os.File, outputBuf *bytes.Buffer, done chan struct{}, cols, rows int) *CapturedOutput {
	// Spinner state with mutex for synchronization
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Gray color
	var spinnerMu sync.Mutex
	helpLinePrinted := false
	spinnerRunning := true
	gray := "\x1b[38;5;240m"
	reset := "\x1b[0m"
	helpText := gray + " Running • Ctrl+C: abort" + reset

	// Channel to stop the spinner
	spinnerStop := make(chan struct{})

	// Spinner goroutine - animates independently
	go func() {
		ticker := time.NewTicker(spinnerTickInterval)
		defer ticker.Stop()
		for {
			select {
			case <-spinnerStop:
				return
			case <-ticker.C:
				spinnerMu.Lock()
				if helpLinePrinted && spinnerRunning {
					s, _ = s.Update(spinner.TickMsg{})
					_, _ = os.Stdout.WriteString("\r" + s.View() + helpText)
				}
				spinnerMu.Unlock()
			}
		}
	}()

	// Start output copying IMMEDIATELY after pty.Start() to avoid delay
	// (terminal setup operations below can be slow on capable terminals)
	go func() {
		buf := make([]byte, ptyReadBufferSize)
		for {
			n, readErr := ptmx.Read(buf)
			if n > 0 {
				spinnerMu.Lock()
				if helpLinePrinted {
					// Clear help line before printing new output
					_, _ = os.Stdout.WriteString("\r\x1b[2K")
					helpLinePrinted = false
				}
				_, _ = os.Stdout.Write(buf[:n])
				outputBuf.Write(buf[:n])

				// Print help line after output
				s, _ = s.Update(spinner.TickMsg{})
				_, _ = os.Stdout.WriteString("\r\n" + s.View() + helpText)
				helpLinePrinted = true
				spinnerMu.Unlock()
			}
			if readErr != nil {
				break
			}
		}
		// Clear help line when done
		spinnerMu.Lock()
		spinnerRunning = false
		if helpLinePrinted {
			_, _ = os.Stdout.WriteString("\r\x1b[2K")
		}
		spinnerMu.Unlock()
		close(spinnerStop)
		close(done)
	}()

	// Set PTY size
	_ = pty.Setsize(ptmx, &pty.Winsize{Rows: uint16(rows - 1), Cols: uint16(cols)})

	// Handle window resize
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGWINCH)
	go func() {
		for range sigCh {
			if w, h, sizeErr := term.GetSize(int(os.Stdout.Fd())); sizeErr == nil {
				_ = pty.Setsize(ptmx, &pty.Winsize{Rows: uint16(h - 1), Cols: uint16(w)})
			}
		}
	}()
	defer signal.Stop(sigCh)

	// Set terminal to raw mode for proper input handling during execution
	oldState, rawErr := term.MakeRaw(int(os.Stdin.Fd()))
	rawModeEnabled := rawErr == nil

	// Top padding and header
	_, _ = os.Stdout.WriteString("\r\n" + gray + "─── Script Output ───" + reset + "\r\n\r\n")

	// Show initial spinner line immediately
	spinnerMu.Lock()
	_, _ = os.Stdout.WriteString(s.View() + helpText)
	helpLinePrinted = true
	spinnerMu.Unlock()

	// Channels for stdin coordination
	scriptDone := make(chan struct{})
	enterPressed := make(chan struct{})

	// Single stdin reader that forwards to PTY during execution,
	// then waits for Enter/Ctrl+C after script finishes
	go func() {
		buf := make([]byte, stdinReadBufferSize)
		for {
			n, readErr := os.Stdin.Read(buf)
			if readErr != nil {
				break
			}
			if n > 0 {
				// Check if script is still running
				select {
				case <-scriptDone:
					// Script finished - check for Enter or Ctrl+C
					for i := 0; i < n; i++ {
						if buf[i] == '\r' || buf[i] == '\n' {
							close(enterPressed)
							return
						}
						if buf[i] == ctrlC { // Ctrl+C - clear help line, restore terminal and exit
							_, _ = os.Stdout.WriteString("\r\x1b[2K\x1b[A\x1b[2K") // Clear line, move up, clear that line too
							if rawModeEnabled {
								_ = term.Restore(int(os.Stdin.Fd()), oldState)
							}
							os.Exit(0)
						}
					}
				default:
					// Script running - forward to PTY
					if _, writeErr := ptmx.Write(buf[:n]); writeErr != nil {
						return
					}
				}
			}
		}
	}()

	// Track execution time
	startTime := time.Now()

	// Wait for command to finish
	cmdErr := cmd.Wait()
	duration := time.Since(startTime)

	// Wait for output goroutine to finish
	<-done

	// Extract exit code
	exitCode := 0
	if cmdErr != nil {
		if exitErr, ok := cmdErr.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}

	// Close PTY
	_ = ptmx.Close()

	// Show completion help line with padding (still in raw mode)
	_, _ = os.Stdout.WriteString("\r\n" + helpLine(false, true))

	// Signal that script is done - stdin reader will now wait for Enter
	close(scriptDone)

	// Wait for Enter key from the stdin reader goroutine
	<-enterPressed

	// Clear the help line and the padding line above it
	_, _ = os.Stdout.WriteString("\r\x1b[2K\x1b[A\x1b[2K")

	// Restore terminal from raw mode
	if rawModeEnabled {
		if restoreErr := term.Restore(int(os.Stdin.Fd()), oldState); restoreErr != nil {
			// Log the error but don't panic - user can reset terminal with 'reset' command
			log.Warn().Err(restoreErr).Msg("Failed to restore terminal state")
		}
	}

	// Reset terminal state explicitly and allow it to settle before Tea starts
	_, _ = os.Stdout.WriteString("\x1b[0m") // Reset all attributes
	_ = os.Stdout.Sync()
	time.Sleep(terminalSettleTimeout)

	// Clean output for history
	rawOutput := outputBuf.String()

	// For failed commands (non-zero exit), preserve raw output to show error messages
	// For successful commands, clean output for better display
	outputToReturn := cleanOutput(rawOutput)
	if exitCode != 0 {
		outputToReturn = rawOutput
	}

	return &CapturedOutput{
		Stdout:   outputToReturn,
		ExitCode: exitCode,
		Duration: duration,
	}
}

// cleanOutput removes control sequences and normalizes the output for storage.
func cleanOutput(raw string) string {
	// Remove OSC sequences and cursor position reports
	cleaned := controlSeqRegex.ReplaceAllString(raw, "")

	// Remove CSI sequences (cursor movement, colors, etc.)
	cleaned = csiRegex.ReplaceAllString(cleaned, "")

	// Remove carriage returns followed by content (overwritten text)
	lines := strings.Split(cleaned, "\n")
	var result []string
	for _, line := range lines {
		// Handle carriage returns - keep only the last segment
		// But if \r is at the end (PTY line ending), just trim it
		if idx := strings.LastIndex(line, "\r"); idx != -1 {
			if idx < len(line)-1 {
				// There's content after \r, keep only that (overwritten text)
				line = line[idx+1:]
			} else {
				// \r is at the end, just remove it
				line = line[:idx]
			}
		}
		// Trim trailing whitespace
		line = strings.TrimRight(line, " \t\r")
		result = append(result, line)
	}

	// Remove trailing empty lines
	for len(result) > 0 && strings.TrimSpace(result[len(result)-1]) == "" {
		result = result[:len(result)-1]
	}

	return strings.Join(result, "\n")
}
