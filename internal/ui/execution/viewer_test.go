package execution

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/Netflix/go-expect"
	pseudotty "github.com/creack/pty"
	"github.com/hinshun/vt10x"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_cleanOutput_Basic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple text", "hello world", "hello world"},
		{"text with trailing newlines", "hello\n\n\n", "hello"},
		{"text with carriage return at end", "hello\r\n", "hello"},
		{"empty input", "", ""},
		{"only whitespace", "   \n\t\n  ", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanOutput(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_cleanOutput_ControlSequences(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"text with OSC sequence", "hello\x1b]0;title\x07world", "helloworld"},
		{"text with CSI color codes", "\x1b[32mgreen\x1b[0m", "green"},
		{"text with cursor movement", "line1\x1b[2Aline2", "line1line2"},
		{"overwritten text with carriage return", "old text\rnew", "new"},
		{"multiple lines with trailing whitespace", "line1   \nline2\t\n", "line1\nline2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanOutput(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_helpLine(t *testing.T) {
	tests := []struct {
		name          string
		running       bool
		fromAssistant bool
		contains      string
	}{
		{"running shows abort message", true, false, "Ctrl+C: abort"},
		{"assistant shows back to assistant", false, true, "back to assistant"},
		{"default shows quit message", false, false, "Enter/Ctrl+C: quit"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helpLine(tt.running, tt.fromAssistant)
			assert.Contains(t, result, tt.contains)
		})
	}
}

func Test_RunWithViewer_Direct(t *testing.T) {
	// Test direct execution (fromAssistant=false) - should not wait for Enter
	pty, tty, err := pseudotty.Open()
	require.NoError(t, err)
	defer func() { _ = pty.Close() }()
	defer func() { _ = tty.Close() }()

	term := vt10x.New(vt10x.WithWriter(tty))
	c, err := expect.NewConsole(expect.WithStdin(pty), expect.WithStdout(term), expect.WithCloser(pty, tty))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	done := make(chan *CapturedOutput, 1)
	go func() {
		oldStdin, oldStdout := os.Stdin, os.Stdout
		os.Stdin = c.Tty()
		os.Stdout = c.Tty()
		defer func() {
			os.Stdin = oldStdin
			os.Stdout = oldStdout
		}()

		cmd := exec.Command("/bin/sh", "-c", "echo 'direct test'")
		result := RunWithViewer(cmd, false) // fromAssistant=false
		done <- result
	}()

	// Wait for result - direct execution should complete without Enter
	select {
	case result := <-done:
		assert.Equal(t, 0, result.ExitCode)
		assert.Contains(t, result.Stdout, "direct test")
	case <-time.After(5 * time.Second):
		t.Fatal("direct execution did not complete in time")
	}
}

func Test_RunWithViewer_Assistant(t *testing.T) {
	// Test assistant execution (fromAssistant=true) - should wait for Enter
	pty, tty, err := pseudotty.Open()
	require.NoError(t, err)
	defer func() { _ = pty.Close() }()
	defer func() { _ = tty.Close() }()

	term := vt10x.New(vt10x.WithWriter(tty))
	c, err := expect.NewConsole(expect.WithStdin(pty), expect.WithStdout(term), expect.WithCloser(pty, tty))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	done := make(chan *CapturedOutput, 1)
	go func() {
		oldStdin, oldStdout := os.Stdin, os.Stdout
		os.Stdin = c.Tty()
		os.Stdout = c.Tty()
		defer func() {
			os.Stdin = oldStdin
			os.Stdout = oldStdout
		}()

		cmd := exec.Command("/bin/sh", "-c", "echo 'assistant test'")
		result := RunWithViewer(cmd, true) // fromAssistant=true
		done <- result
	}()

	// Give time for script to execute
	time.Sleep(200 * time.Millisecond)

	// Send Enter to continue after script finishes
	_, _ = c.Send("\n")

	// Wait for result
	select {
	case result := <-done:
		assert.Equal(t, 0, result.ExitCode)
		assert.Contains(t, result.Stdout, "assistant test")
	case <-time.After(5 * time.Second):
		t.Fatal("assistant execution did not complete in time")
	}
}

func Test_RunWithViewer_ExitCode(t *testing.T) {
	pty, tty, err := pseudotty.Open()
	require.NoError(t, err)
	defer func() { _ = pty.Close() }()
	defer func() { _ = tty.Close() }()

	term := vt10x.New(vt10x.WithWriter(tty))
	c, err := expect.NewConsole(expect.WithStdin(pty), expect.WithStdout(term), expect.WithCloser(pty, tty))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	done := make(chan *CapturedOutput, 1)
	go func() {
		oldStdin, oldStdout := os.Stdin, os.Stdout
		os.Stdin = c.Tty()
		os.Stdout = c.Tty()
		defer func() {
			os.Stdin = oldStdin
			os.Stdout = oldStdout
		}()

		cmd := exec.Command("/bin/sh", "-c", "exit 42")
		result := RunWithViewer(cmd, false)
		done <- result
	}()

	select {
	case result := <-done:
		assert.Equal(t, 42, result.ExitCode)
	case <-time.After(5 * time.Second):
		t.Fatal("execution did not complete in time")
	}
}

func Test_RunWithViewer_Duration(t *testing.T) {
	pty, tty, err := pseudotty.Open()
	require.NoError(t, err)
	defer func() { _ = pty.Close() }()
	defer func() { _ = tty.Close() }()

	term := vt10x.New(vt10x.WithWriter(tty))
	c, err := expect.NewConsole(expect.WithStdin(pty), expect.WithStdout(term), expect.WithCloser(pty, tty))
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	done := make(chan *CapturedOutput, 1)
	go func() {
		oldStdin, oldStdout := os.Stdin, os.Stdout
		os.Stdin = c.Tty()
		os.Stdout = c.Tty()
		defer func() {
			os.Stdin = oldStdin
			os.Stdout = oldStdout
		}()

		cmd := exec.Command("/bin/sh", "-c", "sleep 0.1")
		result := RunWithViewer(cmd, false)
		done <- result
	}()

	select {
	case result := <-done:
		assert.GreaterOrEqual(t, result.Duration, 100*time.Millisecond)
	case <-time.After(5 * time.Second):
		t.Fatal("execution did not complete in time")
	}
}
