package chat

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/phuslu/log"
	"github.com/spf13/afero"

	"github.com/lemoony/snipkit/internal/assistant"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui/style"
)

// UIMode represents the current mode of the unified chat interface.
type UIMode int

const (
	UIModeInput         UIMode = iota // User typing prompt
	UIModeGenerating                  // Async generation with spinner
	UIModeScriptReady                 // Script ready, awaiting user action
	UIModePostExecution               // After execution, show save options
)

const (
	placeholderGenerating = "generating"
	textInitializing      = "\n  Initializing..."
	maxScriptPreviewLen   = 100
	modalPaddingH         = 4
	responsivePaddingH    = 1   // Horizontal padding when terminal is wide enough
	minWidthForPadding    = 100 // Minimum terminal width to apply padding
)

func (m UIMode) String() string {
	return [...]string{
		"Input",
		"Generating",
		"ScriptReady",
		"PostExecution",
	}[m]
}

// UnifiedConfig contains configuration for the unified chat interface.
type UnifiedConfig struct {
	History    []HistoryEntry
	Generating bool
	ScriptChan chan interface{}
	Parameters []model.Parameter
}

// actionBarOption represents a single option in the action bar.
type actionBarOption struct {
	label    string        // Display label (e.g., "Execute")
	shortcut string        // Keyboard shortcut (e.g., "E")
	action   PreviewAction // Action to perform
}

// unifiedChatModel combines chat and preview functionality into a single model.
type unifiedChatModel struct {
	// Common fields
	messages      []ChatMessage
	viewport      viewport.Model
	width, height int
	ready         bool
	currentMode   UIMode

	// Input mode (from chatModel)
	input textinput.Model

	// Action/menu state (from previewModel)
	selectedOption int

	// Generation state (from previewModel)
	generating      bool
	spinnerFrame    int
	scriptChan      chan interface{}
	generatedScript interface{}

	// Modal state (from previewModel)
	modalState      modalState
	paramModal      *parameterModal
	saveModal       *saveModal
	parameters      []model.Parameter
	parameterValues []string

	// Save and execution state (from previewModel)
	saveFilename       string
	saveSnippetName    string
	hasExecutionOutput bool

	// Return values
	action       PreviewAction
	latestPrompt string
	quitting     bool

	styler style.Style
}

// newUnifiedChatModel creates a new unified chat model with the given configuration.
func newUnifiedChatModel(config UnifiedConfig, styler style.Style) *unifiedChatModel {
	messages := buildMessagesFromHistory(config.History)

	// Check if history has execution output
	hasExecutionOutput := false
	for _, entry := range config.History {
		if entry.ExecutionOutput != "" {
			hasExecutionOutput = true
			break
		}
	}

	// Determine initial mode based on config
	initialMode := UIModeInput
	var generatedScriptContent string
	if config.Generating {
		initialMode = UIModeGenerating
		// Add generating placeholder message
		messages = append(messages, ChatMessage{
			Type:    MessageTypeScript,
			Content: placeholderGenerating,
		})
	} else if len(config.History) > 0 {
		lastEntry := config.History[len(config.History)-1]
		if lastEntry.GeneratedScript != "" && lastEntry.ExecutionOutput == "" {
			// Script generated but not executed - show script ready mode
			initialMode = UIModeScriptReady
			generatedScriptContent = lastEntry.GeneratedScript
		} else if lastEntry.ExecutionOutput != "" {
			// Has execution output - show post-execution mode
			initialMode = UIModePostExecution
			generatedScriptContent = lastEntry.GeneratedScript
		}
	}

	m := &unifiedChatModel{
		messages:           messages,
		styler:             styler,
		currentMode:        initialMode,
		selectedOption:     0, // Default to first menu option
		generating:         config.Generating,
		scriptChan:         config.ScriptChan,
		parameters:         config.Parameters,
		modalState:         modalNone,
		hasExecutionOutput: hasExecutionOutput,
		action:             PreviewActionCancel, // Default action
	}

	// Set generatedScript from history if we have it
	// We need to reconstruct a ParsedScript object from the content
	if generatedScriptContent != "" {
		// Create a ParsedScript object with the contents from history
		m.generatedScript = assistant.ParsedScript{
			Contents: generatedScriptContent,
			Filename: "",
			Title:    "",
		}
	}

	// Setup input mode if needed
	if initialMode == UIModeInput {
		m.setupInput()
	}

	return m
}

// setupInput initializes the text input field for prompt entry.
func (m *unifiedChatModel) setupInput() {
	m.input = setupTextInput(m.styler, len(m.messages) > 0)
}

// Init initializes the unified chat model and returns the initial command.
func (m *unifiedChatModel) Init() tea.Cmd {
	log.Trace().
		Str("mode", m.currentMode.String()).
		Int("parameters_count", len(m.parameters)).
		Int("parameter_values_count", len(m.parameterValues)).
		Bool("has_script", m.generatedScript != nil).
		Msg("Initializing unified chat model")

	var cmds []tea.Cmd

	// Always request window size to ensure viewport is properly initialized
	// (especially important when returning from script execution)
	cmds = append(cmds, tea.WindowSize())

	// Handle generation mode
	if m.currentMode == UIModeGenerating {
		cmds = append(cmds, tick())
		if m.scriptChan != nil {
			cmds = append(cmds, waitForScript(m.scriptChan))
		}
	}

	// Handle input mode
	if m.currentMode == UIModeInput {
		cmds = append(cmds, textinput.Blink)
	}

	return tea.Batch(cmds...)
}

// transitionToMode transitions the model to a new UI mode.
func (m *unifiedChatModel) transitionToMode(newMode UIMode) tea.Cmd {
	log.Trace().
		Str("from_mode", m.currentMode.String()).
		Str("to_mode", newMode.String()).
		Msg("Transitioning mode")
	m.currentMode = newMode

	// Recalculate viewport height since bottom bar height may have changed
	if m.ready {
		titleHeight := lipgloss.Height(m.styler.Title("SnipKit Assistant"))
		bottomBarHeight := m.getBottomBarHeight()
		margins := 0
		viewportHeight := m.height - titleHeight - bottomBarHeight - margins

		if viewportHeight > 0 {
			m.viewport.Height = viewportHeight
		}
	}

	switch newMode {
	case UIModeInput:
		return m.setupInputMode()
	case UIModeGenerating:
		return m.setupGeneratingMode()
	case UIModeScriptReady:
		return m.setupScriptReadyMode()
	case UIModePostExecution:
		return m.setupPostExecutionMode()
	}

	return nil
}

// setupInputMode transitions to input mode.
func (m *unifiedChatModel) setupInputMode() tea.Cmd {
	m.setupInput()

	// Set input width (same logic as in handleWindowSize)
	inputWidth := m.width - len(m.input.Prompt) - 2
	if inputWidth < 1 {
		inputWidth = 1
	}
	m.input.Width = inputWidth

	return textinput.Blink
}

// setupGeneratingMode transitions to generating mode.
func (m *unifiedChatModel) setupGeneratingMode() tea.Cmd {
	m.generating = true
	m.spinnerFrame = 0

	// Add generating placeholder message
	m.messages = append(m.messages, ChatMessage{
		Type:    MessageTypeScript,
		Content: placeholderGenerating,
	})

	if m.ready {
		m.viewport.SetContent(m.renderMessagesWithSpinner())
		m.viewport.GotoBottom()
	}

	// Start spinner and await script
	var cmds []tea.Cmd
	cmds = append(cmds, tick())
	if m.scriptChan != nil {
		cmds = append(cmds, waitForScript(m.scriptChan))
	}

	return tea.Batch(cmds...)
}

// setupScriptReadyMode transitions to script ready mode.
func (m *unifiedChatModel) setupScriptReadyMode() tea.Cmd {
	m.generating = false
	m.selectedOption = 0 // Default to Execute
	return nil
}

// setupPostExecutionMode transitions to post-execution mode.
func (m *unifiedChatModel) setupPostExecutionMode() tea.Cmd {
	m.hasExecutionOutput = true
	m.selectedOption = 0 // Default to Execute again
	return nil
}

// Update handles incoming messages and updates the model state.
func (m *unifiedChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle modal overlays first (highest priority)
	if m.modalState != modalNone {
		return m.handleModalUpdate(msg)
	}

	// Handle global messages (window size, mouse, quit, etc.)
	if handled, model, cmd := m.handleGlobalMessages(msg); handled {
		return model, cmd
	}

	// Mode-specific handling
	// Log mode for key messages
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		log.Trace().
			Str("mode", m.currentMode.String()).
			Str("key", keyMsg.String()).
			Msg("Processing key in mode")
	}

	switch m.currentMode {
	case UIModeInput:
		return m.handleInputMode(msg)
	case UIModeGenerating:
		return m.handleGeneratingMode(msg)
	case UIModeScriptReady:
		return m.handleScriptReadyMode(msg)
	case UIModePostExecution:
		return m.handlePostExecutionMode(msg)
	}

	return m, nil
}

// handleGlobalMessages handles messages that apply globally regardless of mode.
// Returns (handled, model, cmd) where handled indicates if the message was processed.
func (m *unifiedChatModel) handleGlobalMessages(msg tea.Msg) (bool, tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Ctrl+C always quits
		if msg.Type == tea.KeyCtrlC {
			m.quitting = true
			m.action = PreviewActionCancel
			return true, m, tea.Quit
		}

		// PgUp/PgDown always goes to viewport for scrolling
		if msg.Type == tea.KeyPgUp || msg.Type == tea.KeyPgDown {
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return true, m, cmd
		}

	case tea.MouseMsg:
		// Handle mouse events for viewport scrolling
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return true, m, cmd

	case tea.WindowSizeMsg:
		m.handleWindowSize(msg)
		return true, m, nil

	case tickMsg:
		if m.generating && m.currentMode == UIModeGenerating {
			m.spinnerFrame++
			if m.ready {
				m.viewport.SetContent(m.renderMessagesWithSpinner())
			}
			return true, m, tick()
		}

	case scriptReadyMsg:
		model, cmd := m.handleScriptReady(msg)
		return true, model, cmd
	}

	return false, m, nil
}

// handleModalUpdate handles updates when a modal is active.
func (m *unifiedChatModel) handleModalUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.modalState {
	case modalParameters:
		return m.handleParameterModal(msg)
	case modalSave:
		return m.handleSaveModal(msg)
	case modalExecuting:
		// Executing overlay is informational only, return immediately
		return m, nil
	}
	return m, nil
}

// handleParameterModal handles updates for the parameter collection modal.
func (m *unifiedChatModel) handleParameterModal(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.paramModal, cmd = m.paramModal.Update(msg)

	if m.paramModal.IsSubmitted() {
		m.parameterValues = m.paramModal.GetValues()
		m.modalState = modalNone
		m.action = PreviewActionExecute
		m.quitting = true
		return m, tea.Quit // Return to app for execution
	}

	if m.paramModal.IsCanceled() {
		m.modalState = modalNone
		// Stay in current mode
		return m, nil
	}

	return m, cmd
}

// handleSaveModal handles updates for the save modal.
func (m *unifiedChatModel) handleSaveModal(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.saveModal, cmd = m.saveModal.Update(msg)

	if m.saveModal.IsSubmitted() {
		m.saveFilename = m.saveModal.GetFilename()
		m.saveSnippetName = m.saveModal.GetSnippetName()
		m.modalState = modalNone
		m.action = PreviewActionCancel // Save and exit
		m.quitting = true
		return m, tea.Quit
	}

	if m.saveModal.IsCanceled() {
		m.modalState = modalNone
		// Return to post-execution mode
		return m, nil
	}

	return m, cmd
}

// handleWindowSize handles terminal window resize events.
func (m *unifiedChatModel) handleWindowSize(msg tea.WindowSizeMsg) {
	m.width = msg.Width
	m.height = msg.Height

	// Calculate effective content width accounting for responsive padding
	contentWidth := msg.Width
	if msg.Width >= minWidthForPadding {
		contentWidth = msg.Width - (responsivePaddingH * 2)
	}

	// Calculate dimensions
	titleHeight := lipgloss.Height(m.styler.Title("SnipKit Assistant"))
	bottomBarHeight := m.getBottomBarHeight() // Dynamically calculated based on current mode
	margins := 0

	viewportHeight := msg.Height - titleHeight - bottomBarHeight - margins
	if viewportHeight < 1 {
		viewportHeight = 1
	}

	if !m.ready {
		// First time setup
		m.viewport = viewport.New(contentWidth, viewportHeight)
		m.viewport.YPosition = 0

		if m.generating {
			m.viewport.SetContent(m.renderMessagesWithSpinner())
		} else {
			m.viewport.SetContent(renderMessages(m.messages, m.styler, contentWidth))
		}
		m.viewport.GotoBottom()
		m.ready = true
	} else {
		// Resize existing viewport
		m.viewport.Width = contentWidth
		m.viewport.Height = viewportHeight

		if m.generating {
			m.viewport.SetContent(m.renderMessagesWithSpinner())
		} else {
			m.viewport.SetContent(renderMessages(m.messages, m.styler, contentWidth))
		}
	}

	// Update input width if in input mode
	if m.currentMode == UIModeInput {
		inputWidth := contentWidth - len(m.input.Prompt) - 2
		if inputWidth < 1 {
			inputWidth = 1
		}
		m.input.Width = inputWidth
	}
}

// handleScriptReady handles the script ready message from async generation.
func (m *unifiedChatModel) handleScriptReady(msg scriptReadyMsg) (tea.Model, tea.Cmd) {
	m.generatedScript = msg.script
	m.generating = false

	// Extract script content and update the last message
	var scriptContent string
	if len(m.messages) > 0 && m.messages[len(m.messages)-1].Content == placeholderGenerating {
		// Extract Contents field using reflection
		v := reflect.ValueOf(msg.script)
		if v.Kind() == reflect.Struct {
			contentsField := v.FieldByName("Contents")
			if contentsField.IsValid() && contentsField.Kind() == reflect.String {
				scriptContent = contentsField.String()
				m.messages[len(m.messages)-1].Content = scriptContent
			}
		}
	}

	// Extract parameters from the generated script
	if scriptContent != "" {
		snippet := assistant.PrepareSnippet([]byte(scriptContent), assistant.ParsedScript{
			Contents: scriptContent,
		})
		m.parameters = snippet.GetParameters()
		log.Trace().
			Int("parameters_count", len(m.parameters)).
			Str("script_preview", scriptContent[:min(maxScriptPreviewLen, len(scriptContent))]).
			Msg("Extracted parameters from generated script")
	} else {
		log.Warn().Msg("No script content to extract parameters from")
	}

	// Update viewport
	if m.ready {
		m.viewport.SetContent(renderMessages(m.messages, m.styler, m.width))
		m.viewport.GotoBottom()
	}

	// Transition to script ready mode
	cmd := m.transitionToMode(UIModeScriptReady)
	return m, cmd
}

// handleInputMode handles input mode updates.
func (m *unifiedChatModel) handleInputMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyEsc:
			m.quitting = true
			m.action = PreviewActionCancel
			return m, tea.Quit

		case tea.KeyEnter:
			m.latestPrompt = m.input.Value()
			if m.latestPrompt != "" {
				m.action = PreviewActionRevise // User entered a new prompt
				m.quitting = true
				return m, tea.Quit
			}
			// Empty input, stay in mode
			return m, nil

		default:
			// Pass to input
			m.input, cmd = m.input.Update(keyMsg)
			return m, cmd
		}
	}

	return m, cmd
}

// handleGeneratingMode handles generating mode updates.
func (m *unifiedChatModel) handleGeneratingMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.Type == tea.KeyEsc {
			// Cancel generation
			m.quitting = true
			m.action = PreviewActionCancel
			return m, tea.Quit
		}
	}

	return m, nil
}

// handleScriptReadyMode handles script ready mode updates.
func (m *unifiedChatModel) handleScriptReadyMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		shouldExecute, action, newSelectedOption := m.handleActionBarInput(
			keyMsg,
			m.getScriptReadyOptions(),
		)

		m.selectedOption = newSelectedOption

		if shouldExecute {
			switch action {
			case PreviewActionExecute:
				return m.handleExecuteAction()
			case PreviewActionEdit:
				m.action = PreviewActionEdit
				m.quitting = true
				return m, tea.Quit
			case PreviewActionRevise:
				cmd := m.transitionToMode(UIModeInput)
				return m, cmd
			case PreviewActionCancel:
				m.action = PreviewActionCancel
				m.quitting = true
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

// handleExecuteAction initiates script execution.
func (m *unifiedChatModel) handleExecuteAction() (tea.Model, tea.Cmd) {
	// If parameters exist and not yet collected, show parameter modal
	if len(m.parameters) > 0 && len(m.parameterValues) == 0 {
		log.Trace().Msg("Showing parameter modal")
		m.modalState = modalParameters
		m.paramModal = NewParameterModal(m.parameters, m.styler, afero.NewOsFs())
		return m, m.paramModal.Init()
	}

	// No parameters or already collected - execute immediately
	log.Trace().Msg("Proceeding to execution (quitting)")
	m.action = PreviewActionExecute
	m.quitting = true
	return m, tea.Quit
}

// handlePostExecutionMode handles post-execution mode updates.
func (m *unifiedChatModel) handlePostExecutionMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		shouldExecute, action, newSelectedOption := m.handleActionBarInput(
			keyMsg,
			m.getPostExecutionOptions(),
		)

		m.selectedOption = newSelectedOption

		if shouldExecute {
			return m.handlePostExecutionAction(action)
		}
	}

	return m, nil
}

// handlePostExecutionAction handles action selection in post-execution mode.
func (m *unifiedChatModel) handlePostExecutionAction(action PreviewAction) (tea.Model, tea.Cmd) {
	switch action {
	case PreviewActionExecute:
		// Execute again - check if we need parameter modal
		log.Trace().
			Int("parameters_count", len(m.parameters)).
			Int("parameter_values_count", len(m.parameterValues)).
			Msg("Execute again action triggered")

		// If we have parameters and haven't collected values yet, show modal
		if len(m.parameters) > 0 && len(m.parameterValues) == 0 {
			log.Trace().Msg("Showing parameter modal for Execute again")
			m.modalState = modalParameters
			m.paramModal = NewParameterModal(m.parameters, m.styler, afero.NewOsFs())
			return m, m.paramModal.Init()
		}
		// Otherwise proceed to execution
		log.Trace().Msg("Execute again: proceeding to execution (quitting)")
		m.action = PreviewActionExecute
		m.quitting = true
		return m, tea.Quit

	case PreviewActionRevise:
		// Transition to input mode for new question
		cmd := m.transitionToMode(UIModeInput)
		return m, cmd

	case PreviewActionCancel:
		// This is "Save & Exit" in post-execution mode
		// Show save modal
		m.modalState = modalSave

		// Extract proposed save values from generated script
		proposedFilename := ""
		proposedSnippetName := ""
		if m.generatedScript != nil {
			v := reflect.ValueOf(m.generatedScript)
			if v.Kind() == reflect.Struct {
				if f := v.FieldByName("Filename"); f.IsValid() && f.Kind() == reflect.String {
					proposedFilename = f.String()
				}
				if f := v.FieldByName("Title"); f.IsValid() && f.Kind() == reflect.String {
					proposedSnippetName = f.String()
				}
			}
		}

		m.saveModal = NewSaveModal(proposedFilename, proposedSnippetName, m.styler, afero.NewOsFs())
		return m, m.saveModal.Init()

	case PreviewActionExitNoSave:
		m.action = PreviewActionExitNoSave
		m.quitting = true
		return m, tea.Quit
	}

	return m, nil
}

// renderMessagesWithSpinner renders messages with spinner animation for generating state.
func (m *unifiedChatModel) renderMessagesWithSpinner() string {
	// Create a copy of messages with the last one showing spinner
	messagesForRender := make([]ChatMessage, len(m.messages))
	copy(messagesForRender, m.messages)

	if len(messagesForRender) > 0 && messagesForRender[len(messagesForRender)-1].Content == placeholderGenerating {
		spinnerFrame := spinnerFrames[m.spinnerFrame%len(spinnerFrames)]
		messagesForRender[len(messagesForRender)-1].Content = fmt.Sprintf("%s Generating script...", spinnerFrame)
	}

	return renderMessages(messagesForRender, m.styler, m.contentWidth())
}

// View renders the unified chat interface.
func (m *unifiedChatModel) View() string {
	if m.quitting {
		return ""
	}

	if !m.ready {
		return textInitializing
	}

	// Build the base view
	var sections []string

	// Title
	sections = append(sections, m.styler.Title("SnipKit Assistant"))

	// Viewport with message history
	sections = append(sections, m.viewport.View())

	// Bottom bar (mode-dependent)
	sections = append(sections, m.renderBottomBar())

	baseView := lipgloss.JoinVertical(lipgloss.Left, sections...)

	// Apply responsive padding
	if m.width >= minWidthForPadding {
		baseView = lipgloss.NewStyle().
			Padding(0, responsivePaddingH).
			Render(baseView)
	}

	// Overlay modals if active
	if m.modalState != modalNone {
		return m.renderModalOverlay(baseView)
	}

	return baseView
}

// renderBottomBar renders the bottom bar based on current mode.
func (m *unifiedChatModel) renderBottomBar() string {
	switch m.currentMode {
	case UIModeInput:
		return m.renderInputBar()
	case UIModeGenerating:
		return m.renderGeneratingBar()
	case UIModeScriptReady:
		return m.renderScriptReadyBar()
	case UIModePostExecution:
		return m.renderPostExecutionBar()
	}
	return ""
}

// getBottomBarHeight returns the actual height of the bottom bar for the current mode.
func (m *unifiedChatModel) getBottomBarHeight() int {
	// Render the bottom bar and measure its actual height
	bottomBar := m.renderBottomBar()
	return lipgloss.Height(bottomBar)
}

// renderInputBar renders the input bar for input mode.
func (m *unifiedChatModel) renderInputBar() string {
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(m.styler.ActiveColor().Value()).
		Background(m.styler.VerySubduedColor().Value()).
		PaddingTop(1).
		PaddingLeft(2).
		PaddingRight(2).
		Width(m.contentWidth())

	inputRow := inputStyle.Render(m.input.View())

	helpStyle := lipgloss.NewStyle().
		Foreground(m.styler.PlaceholderColor().Value()).
		Padding(0, 2)
	helpText := helpStyle.Render("Enter: submit • Esc: cancel • PgUp/PgDown: scroll")

	// The inputStyle already has vertical padding, so no extra spacing needed
	return lipgloss.JoinVertical(lipgloss.Left, inputRow, helpText)
}

// renderGeneratingBar renders the bar for generating mode.
func (m *unifiedChatModel) renderGeneratingBar() string {
	spinnerFrame := spinnerFrames[m.spinnerFrame%len(spinnerFrames)]
	text := fmt.Sprintf("%s Generating script... [C] Cancel", spinnerFrame)

	style := lipgloss.NewStyle().
		Foreground(m.styler.PlaceholderColor().Value()).
		Padding(0, 2)

	return style.Render(text)
}

// getScriptReadyOptions returns options for script ready mode.
func (m *unifiedChatModel) getScriptReadyOptions() []actionBarOption {
	return []actionBarOption{
		{label: "Execute", shortcut: "E", action: PreviewActionExecute},
		{label: "Open editor", shortcut: "O", action: PreviewActionEdit},
		{label: "Revise", shortcut: "R", action: PreviewActionRevise},
		{label: "Cancel", shortcut: "C", action: PreviewActionCancel},
	}
}

// getPostExecutionOptions returns options for post-execution mode.
func (m *unifiedChatModel) getPostExecutionOptions() []actionBarOption {
	return []actionBarOption{
		{label: "Execute again", shortcut: "E", action: PreviewActionExecute},
		{label: "Revise", shortcut: "R", action: PreviewActionRevise},
		{label: "Save & Exit", shortcut: "S", action: PreviewActionCancel}, // Save
		{label: "Exit & Don't save", shortcut: "X", action: PreviewActionExitNoSave},
	}
}

// renderActionBar renders an action bar with the given options.
func (m *unifiedChatModel) renderActionBar(options []actionBarOption) string {
	var items []string
	for i, opt := range options {
		style := m.getMenuItemStyle(i == m.selectedOption)
		label := fmt.Sprintf("[%s] %s", opt.shortcut, opt.label)
		items = append(items, style.Render(label))
	}

	menu := lipgloss.JoinHorizontal(lipgloss.Top, items...)
	helpText := lipgloss.NewStyle().
		Foreground(m.styler.PlaceholderColor().Value()).
		Render("\n  ←/→ to select • Enter to confirm • Shortcuts available • Esc to cancel")

	style := lipgloss.NewStyle().
		Padding(0, 2)

	return style.Render(fmt.Sprintf("%s%s", menu, helpText))
}

// handleActionBarInput processes input for action bar navigation.
// Returns: (shouldExecute bool, action PreviewAction, newSelectedOption int).
func (m *unifiedChatModel) handleActionBarInput(
	keyMsg tea.KeyMsg,
	options []actionBarOption,
) (bool, PreviewAction, int) {
	numOptions := len(options)

	switch keyMsg.Type {
	case tea.KeyEsc:
		return true, PreviewActionCancel, m.selectedOption

	case tea.KeyLeft:
		if m.selectedOption > 0 {
			return false, -1, m.selectedOption - 1
		}
		return false, -1, m.selectedOption

	case tea.KeyRight:
		if m.selectedOption < numOptions-1 {
			return false, -1, m.selectedOption + 1
		}
		return false, -1, m.selectedOption

	case tea.KeyEnter:
		return true, options[m.selectedOption].action, m.selectedOption
	}

	// Check for keyboard shortcuts
	key := strings.ToLower(keyMsg.String())
	for _, opt := range options {
		if strings.ToLower(opt.shortcut) == key {
			return true, opt.action, m.selectedOption
		}
	}

	return false, -1, m.selectedOption
}

// renderScriptReadyBar renders the action bar for script ready mode.
func (m *unifiedChatModel) renderScriptReadyBar() string {
	return m.renderActionBar(m.getScriptReadyOptions())
}

// renderPostExecutionBar renders the action bar for post-execution mode.
func (m *unifiedChatModel) renderPostExecutionBar() string {
	return m.renderActionBar(m.getPostExecutionOptions())
}

// getMenuItemStyle returns the style for a menu item.
func (m *unifiedChatModel) getMenuItemStyle(selected bool) lipgloss.Style {
	if selected {
		return lipgloss.NewStyle().
			Foreground(m.styler.HighlightColor().Value()).
			Background(m.styler.ActiveColor().Value()).
			Bold(true).
			Padding(0, 1)
	}
	return lipgloss.NewStyle().
		Foreground(m.styler.TextColor().Value()).
		Padding(0, 1)
}

// renderModalOverlay renders a modal overlay over the base view.
func (m *unifiedChatModel) renderModalOverlay(baseView string) string {
	switch m.modalState {
	case modalParameters:
		if m.paramModal != nil {
			modal := m.paramModal.View(m.width, m.height)
			// Darken background by rendering a semi-transparent overlay
			return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modal)
		}
	case modalSave:
		if m.saveModal != nil {
			modal := m.saveModal.View(m.width, m.height)
			return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modal)
		}
	case modalExecuting:
		executingText := "⚡ Executing script..."
		executingStyle := lipgloss.NewStyle().
			Foreground(m.styler.HighlightColor().Value()).
			Background(m.styler.VerySubduedColor().Value()).
			Bold(true).
			Padding(2, modalPaddingH).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(m.styler.ActiveColor().Value())
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, executingStyle.Render(executingText))
	}
	return baseView
}

// contentWidth returns the effective content width accounting for responsive padding.
func (m *unifiedChatModel) contentWidth() int {
	if m.width >= minWidthForPadding {
		return m.width - (responsivePaddingH * 2)
	}
	return m.width
}
