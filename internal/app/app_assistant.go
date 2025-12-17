package app

import (
	"os"
	"strings"
	"time"

	"github.com/phuslu/log"

	"github.com/lemoony/snipkit/internal/assistant"
	"github.com/lemoony/snipkit/internal/config"
	"github.com/lemoony/snipkit/internal/managers"
	"github.com/lemoony/snipkit/internal/managers/fslibrary"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui"
	"github.com/lemoony/snipkit/internal/ui/assistant/chat"
	"github.com/lemoony/snipkit/internal/ui/picker"
	"github.com/lemoony/snipkit/internal/ui/uimsg"
	"github.com/lemoony/snipkit/internal/utils/sliceutil"
	"github.com/lemoony/snipkit/internal/utils/stringutil"
	"github.com/lemoony/snipkit/internal/utils/tmpdir"
)

func (a *appImpl) GenerateSnippetWithAssistant(demoScriptPath []string, demoQueryDuration time.Duration) {
	asst := a.assistantProviderFunc(
		a.config.Assistant,
		assistant.DemoConfig{ScriptPaths: demoScriptPath, QueryDuration: demoQueryDuration},
	)

	if valid, msg := asst.Initialize(); !valid {
		a.tui.PrintAndExit(msg, -1)
	}

	// Start unified assistant loop with empty history
	history := []chat.HistoryEntry{}
	a.unifiedAssistantLoop(history, asst)
}

func (a *appImpl) saveScript(contents []byte, title, filename string) {
	if manager, ok := sliceutil.FindElement(a.managers, func(manager managers.Manager) bool {
		return manager.Key() == fslibrary.Key
	}); ok {
		manager.SaveAssistantSnippet(title, filename, contents)
	} else {
		panic("File system library not configured as manager. Try running `snipkit manager add`")
	}
}

// handleCancelAction handles the cancel/save action.
func (a *appImpl) handleCancelAction(scriptInterface interface{}, saveFilename, saveSnippetName string) {
	// Check if save data was provided
	if saveFilename != "" || saveSnippetName != "" {
		// User saved from modal
		if scriptInterface != nil {
			if parsed, ok := scriptInterface.(assistant.ParsedScript); ok {
				snippet := assistant.PrepareSnippet([]byte(parsed.Contents), parsed)
				filename := stringutil.StringOrDefault(saveFilename, assistant.RandomScriptFilename())
				title := stringutil.StringOrDefault(saveSnippetName, parsed.Title)

				log.Debug().
					Str("title", title).
					Str("filename", filename).
					Msg("Saving assistant-generated snippet")

				a.saveScript([]byte(snippet.GetContent()), title, filename)
			}
		}
	}
}

// handleReviseAction handles the revise action. Returns (shouldReturn, updatedHistory).
func (a *appImpl) handleReviseAction(history []chat.HistoryEntry, scriptInterface interface{}, latestPrompt string) (bool, []chat.HistoryEntry) {
	if latestPrompt == "" {
		log.Warn().Msg("PreviewActionRevise but no prompt provided")
		return true, history
	}

	// Update last history entry with current script if present
	if scriptInterface != nil {
		if parsed, ok := scriptInterface.(assistant.ParsedScript); ok && len(history) > 0 {
			lastIdx := len(history) - 1
			history[lastIdx].GeneratedScript = parsed.Contents
		}
	}

	// Add new history entry with the new prompt
	history = append(history, chat.HistoryEntry{
		UserPrompt: latestPrompt,
	})

	return false, history
}

// handleExecuteAction handles the execute action and returns updated history.
func (a *appImpl) handleExecuteAction(history []chat.HistoryEntry, scriptInterface interface{}, paramValues []string) []chat.HistoryEntry {
	if scriptInterface == nil {
		log.Error().Msg("Execute action but no script available")
		return a.addExecutionError(history, "Error: No script available to execute")
	}

	parsed, ok := scriptInterface.(assistant.ParsedScript)
	if !ok {
		log.Error().Msgf("Script interface is wrong type: %T", scriptInterface)
		return a.addExecutionError(history, "Error: Invalid script type")
	}

	// Prepare snippet and check for parameters
	snippet := assistant.PrepareSnippet([]byte(parsed.Contents), parsed)
	parameters := snippet.GetParameters()

	if len(parameters) > 0 && len(paramValues) == 0 {
		log.Warn().Msgf("Executing script with %d parameters but no values provided", len(parameters))
	}

	// Execute the snippet
	log.Trace().Msg("About to execute snippet")
	capturedResult := a.executeSnippet(ContextAssistant, false, snippet, paramValues)
	executionTime := time.Now()
	log.Trace().Msg("Snippet execution completed, about to return to chat")

	return a.updateHistoryWithSuccess(history, parsed.Contents, capturedResult, executionTime)
}

// handleEditAction handles the edit action. Returns (shouldContinue, updatedHistory).
func (a *appImpl) handleEditAction(history []chat.HistoryEntry, scriptInterface interface{}, tmpDirSvc tmpdir.TmpDir) (bool, []chat.HistoryEntry) {
	if scriptInterface == nil {
		log.Warn().Msg("Edit action but no script available")
		return false, history
	}

	parsed, ok := scriptInterface.(assistant.ParsedScript)
	if !ok {
		log.Warn().Msg("Script interface is not ParsedScript type")
		return false, history
	}

	fileOk, filePath := tmpDirSvc.CreateTempFile([]byte(parsed.Contents))
	if !fileOk {
		return false, history
	}

	a.tui.OpenEditor(filePath, a.config.Editor)
	//nolint:gosec // ignore potential file inclusion via variable
	updatedContents, err := os.ReadFile(filePath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read edited file")
		return false, history
	}

	snippet := assistant.PrepareSnippet(updatedContents, parsed)
	parameters := snippet.GetParameters()

	var editedParamValues []string
	if len(parameters) > 0 {
		var paramOk bool
		editedParamValues, paramOk = a.tui.ShowParameterForm(parameters, nil, ui.OkButtonExecute)
		if !paramOk {
			return false, history
		}
	}

	// Execute the edited snippet
	capturedResult := a.executeSnippet(ContextAssistant, false, snippet, editedParamValues)

	// Update history with execution results
	if len(history) > 0 {
		lastIdx := len(history) - 1
		executionTime := time.Now()
		history[lastIdx].GeneratedScript = string(updatedContents)
		history[lastIdx].ExecutionOutput = capturedResult.stdout + capturedResult.stderr
		history[lastIdx].ExitCode = &capturedResult.exitCode
		history[lastIdx].Duration = &capturedResult.duration
		history[lastIdx].ExecutionTime = &executionTime
	}

	return true, history
}

// addExecutionError adds an error to the last history entry.
func (a *appImpl) addExecutionError(history []chat.HistoryEntry, errorMsg string) []chat.HistoryEntry {
	if len(history) > 0 {
		lastIdx := len(history) - 1
		history[lastIdx].ExecutionOutput = errorMsg
		exitCode := 1
		history[lastIdx].ExitCode = &exitCode
	}
	return history
}

// updateHistoryWithSuccess updates history with successful execution results.
func (a *appImpl) updateHistoryWithSuccess(history []chat.HistoryEntry, scriptContents string, result *capturedOutput, executionTime time.Time) []chat.HistoryEntry {
	isExecuteAgain := len(history) > 0 && history[len(history)-1].ExecutionOutput != ""

	if isExecuteAgain {
		log.Trace().Msg("Execute again: appending new history entry")
		return append(history, chat.HistoryEntry{
			UserPrompt:      "",
			GeneratedScript: scriptContents,
			ExecutionOutput: result.stdout + result.stderr,
			ExitCode:        &result.exitCode,
			Duration:        &result.duration,
			ExecutionTime:   &executionTime,
		})
	}

	log.Trace().Msg("First execution: updating existing history entry")
	if len(history) > 0 {
		lastIdx := len(history) - 1
		history[lastIdx].GeneratedScript = scriptContents
		history[lastIdx].ExecutionOutput = result.stdout + result.stderr
		history[lastIdx].ExitCode = &result.exitCode
		history[lastIdx].Duration = &result.duration
		history[lastIdx].ExecutionTime = &executionTime
	}
	return history
}

// unifiedAssistantLoop manages the unified chat interaction loop.
func (a *appImpl) unifiedAssistantLoop(history []chat.HistoryEntry, asst assistant.Assistant) {
	tmpDirSvc := tmpdir.New(a.system)
	defer tmpDirSvc.ClearFiles()

	for {
		// Build config based on current history state
		config := a.buildUnifiedConfig(history, asst)

		// Show unified chat (handles all modes internally)
		scriptInterface, paramValues, action, latestPrompt, saveFilename, saveSnippetName := a.tui.ShowUnifiedAssistantChat(config)

		// Handle the action
		switch action {
		case chat.PreviewActionCancel:
			a.handleCancelAction(scriptInterface, saveFilename, saveSnippetName)
			return

		case chat.PreviewActionExitNoSave:
			return

		case chat.PreviewActionRevise:
			if shouldReturn, updatedHistory := a.handleReviseAction(history, scriptInterface, latestPrompt); shouldReturn {
				return
			} else {
				history = updatedHistory
				continue
			}

		case chat.PreviewActionExecute:
			history = a.handleExecuteAction(history, scriptInterface, paramValues)
			continue

		case chat.PreviewActionEdit:
			if shouldContinue, updatedHistory := a.handleEditAction(history, scriptInterface, tmpDirSvc); shouldContinue {
				history = updatedHistory
				continue
			}
			return
		}
	}
}

// buildUnifiedConfig creates a UnifiedConfig based on the current history state.
func (a *appImpl) extractParameters(scriptContent string) []model.Parameter {
	if scriptContent == "" {
		return nil
	}
	snippet := assistant.PrepareSnippet([]byte(scriptContent), assistant.ParsedScript{
		Contents: scriptContent,
	})
	return snippet.GetParameters()
}

func (a *appImpl) buildUnifiedConfig(history []chat.HistoryEntry, asst assistant.Assistant) chat.UnifiedConfig {
	// Case 1: Empty history or new prompt needed - start in input mode
	if len(history) == 0 {
		return chat.UnifiedConfig{
			History:    history,
			Generating: false,
			ScriptChan: nil,
			Parameters: nil,
		}
	}

	lastEntry := history[len(history)-1]

	// Case 2: User entered prompt but no script generated yet - start generation
	if lastEntry.UserPrompt != "" && lastEntry.GeneratedScript == "" {
		// Start async generation
		prompt := lastEntry.UserPrompt
		scriptChan := make(chan interface{}, 1)
		go func() {
			script := asst.Query(prompt)
			scriptChan <- script
		}()

		return chat.UnifiedConfig{
			History:    history,
			Generating: true,
			ScriptChan: scriptChan,
			Parameters: nil,
		}
	}

	// Case 3: Script generated but not executed - show action menu mode
	if lastEntry.GeneratedScript != "" && lastEntry.ExecutionOutput == "" {
		return chat.UnifiedConfig{
			History:    history,
			Generating: false,
			ScriptChan: nil,
			Parameters: a.extractParameters(lastEntry.GeneratedScript),
		}
	}

	// Case 4: Script executed - show post-execution mode
	if lastEntry.ExecutionOutput != "" {
		return chat.UnifiedConfig{
			History:    history,
			Generating: false,
			ScriptChan: nil,
			Parameters: a.extractParameters(lastEntry.GeneratedScript),
		}
	}

	// Default: input mode
	return chat.UnifiedConfig{
		History:    history,
		Generating: false,
		ScriptChan: nil,
		Parameters: nil,
	}
}

func (a *appImpl) EnableAssistant() {
	assistantInstance := a.assistantProviderFunc(a.config.Assistant, assistant.DemoConfig{})
	assistantDescriptions := assistantInstance.AssistantDescriptions(a.config.Assistant)

	listItems := make([]picker.Item, len(assistantDescriptions))
	var selected *picker.Item
	for i := range assistantDescriptions {
		listItems[i] = picker.NewItem(assistantDescriptions[i].Name, assistantDescriptions[i].Description)
		if assistantDescriptions[i].Enabled {
			selected = &listItems[i]
		}
	}

	if selectedIndex, ok := a.tui.ShowPicker("Which AI provider for the assistant do you want to use?", listItems, selected); ok {
		assistantDescription := assistantDescriptions[selectedIndex]
		log.Debug().
			Str("provider", assistantDescription.Name).
			Msg("User selected AI assistant provider")

		// Serialize current config
		oldConfigBytes := config.SerializeToYamlWithComment(config.Wrap(*a.config))
		oldConfigStr := strings.TrimSpace(string(oldConfigBytes))

		// Get new assistant config
		cfg := assistantInstance.AutoConfig(assistantDescription.Key)

		// Create a copy of current config and apply the new assistant to show full diff
		newConfig := *a.config
		newConfig.Assistant = cfg

		// Serialize new config
		newConfigBytes := config.SerializeToYamlWithComment(config.Wrap(newConfig))
		newConfigStr := strings.TrimSpace(string(newConfigBytes))

		// Pass both configs
		confirmed := a.tui.Confirmation(uimsg.ManagerConfigAddConfirm(oldConfigStr, newConfigStr))
		if confirmed {
			a.configService.UpdateAssistantConfig(cfg)
			log.Debug().Str("provider", assistantDescription.Name).Msg("Assistant configuration updated")
		}
		a.tui.Print(uimsg.AssistantUpdateConfigResult(confirmed, a.configService.ConfigFilePath()))
	}
}
