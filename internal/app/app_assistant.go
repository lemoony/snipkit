package app

import (
	"fmt"
	"os"
	"strings"
	"time"

	"emperror.dev/errors"
	"github.com/phuslu/log"

	"github.com/lemoony/snipkit/internal/assistant"
	"github.com/lemoony/snipkit/internal/config"
	"github.com/lemoony/snipkit/internal/managers"
	"github.com/lemoony/snipkit/internal/managers/fslibrary"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui"
	"github.com/lemoony/snipkit/internal/ui/assistant/chat"
	"github.com/lemoony/snipkit/internal/ui/assistant/wizard"
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

	if ok, text := a.tui.ShowAssistantPrompt([]chat.HistoryEntry{}); ok {
		history := []chat.HistoryEntry{{UserPrompt: text}}
		a.generateAndHandleScript(text, history, asst)
	}
}

func (a *appImpl) generateAndHandleScript(prompt string, history []chat.HistoryEntry, asst assistant.Assistant) {
	// Create a function that generates the script
	generateScript := func() interface{} {
		return asst.Query(prompt)
	}

	// Show preview with async generation
	scriptInterface, action := a.tui.ShowAssistantScriptPreviewWithGeneration(history, generateScript)

	if action == chat.PreviewActionCancel {
		return
	}

	// Convert back to ParsedScript
	script, ok := scriptInterface.(assistant.ParsedScript)
	if !ok {
		return
	}

	// Handle the generated script with the chosen action
	a.handleGeneratedScriptWithAction(script, action, history, asst)
}

func (a *appImpl) handleGeneratedScriptWithAction(parsed assistant.ParsedScript, action chat.PreviewAction, history []chat.HistoryEntry, asst assistant.Assistant) {
	tmpDirSvc := tmpdir.New(a.system)
	defer tmpDirSvc.ClearFiles()

	switch action {
	case chat.PreviewActionCancel:
		// User canceled the preview
		return

	case chat.PreviewActionRevise:
		// User wants to revise with a new prompt without executing
		// First, update the last history entry with the generated script
		if len(history) > 0 {
			lastIdx := len(history) - 1
			history[lastIdx].GeneratedScript = parsed.Contents
		}

		if ok, newPrompt := a.tui.ShowAssistantPrompt(history); ok {
			// Add new entry for revision
			history = append(history, chat.HistoryEntry{
				UserPrompt: newPrompt,
			})
			a.generateAndHandleScript(newPrompt, history, asst)
		}
		return

	case chat.PreviewActionExecute:
		// Execute directly without opening editor
		snippet := assistant.PrepareSnippet([]byte(parsed.Contents), parsed)
		if parameters := snippet.GetParameters(); len(parameters) > 0 {
			parameterValues, paramOk := a.tui.ShowParameterForm(parameters, nil, ui.OkButtonExecute)
			if paramOk {
				a.executeAndHandleSnippet(snippet, parameterValues, history, asst, parsed)
			}
		} else {
			a.executeAndHandleSnippet(snippet, nil, history, asst, parsed)
		}

	case chat.PreviewActionEdit:
		// Open in editor first, then execute
		if fileOk, filePath := tmpDirSvc.CreateTempFile([]byte(parsed.Contents)); fileOk {
			a.tui.OpenEditor(filePath, a.config.Editor)
			//nolint:gosec // ignore potential file inclusion via variable
			updatedContents, err := os.ReadFile(filePath)
			if err != nil {
				panic(errors.Wrapf(err, "failed to read temporary file"))
			}

			snippet := assistant.PrepareSnippet(updatedContents, parsed)
			if parameters := snippet.GetParameters(); len(parameters) > 0 {
				parameterValues, paramOk := a.tui.ShowParameterForm(parameters, nil, ui.OkButtonExecute)
				if paramOk {
					a.executeAndHandleSnippet(snippet, parameterValues, history, asst, parsed)
				}
			} else {
				a.executeAndHandleSnippet(snippet, nil, history, asst, parsed)
			}
		}
	}
}

func (a *appImpl) executeAndHandleSnippet(snippet model.Snippet, parameterValues []string, history []chat.HistoryEntry, asst assistant.Assistant, script assistant.ParsedScript) {
	executed, capturedResult := a.executeSnippet(len(parameterValues) == 0, false, snippet, parameterValues)
	if executed {
		// Update last history entry with script and output
		if len(history) > 0 {
			lastIdx := len(history) - 1
			history[lastIdx].GeneratedScript = script.Contents
			history[lastIdx].ExecutionOutput = capturedResult.stdout + capturedResult.stderr
		}

		wizardOk, result := a.tui.ShowAssistantWizard(wizard.Config{
			ShowSaveOption:      a.config.Assistant.SaveMode != assistant.SaveModeNever,
			ProposedFilename:    script.Filename,
			ProposedSnippetName: script.Title,
		})
		if wizardOk {
			switch result.SelectedOption {
			case wizard.OptionTryAgain:
				log.Debug().Msg("User requested to try again with assistant")
				if ok2, prompt2 := a.tui.ShowAssistantPrompt(history); ok2 {
					// Add new entry for retry
					history = append(history, chat.HistoryEntry{
						UserPrompt: prompt2,
					})
					newPrompt := fmt.Sprintf("The result of the command was: %s\n%s\n\n%s", capturedResult.stdout, capturedResult.stderr, prompt2)
					a.generateSnippetWithAdditionalPrompt(newPrompt, history, asst)
				}
			case wizard.OptionSaveExit:
				log.Debug().
					Str("title", result.SnippetTitle).
					Str("filename", stringutil.StringOrDefault(result.Filename, assistant.RandomScriptFilename())).
					Msg("Saving assistant-generated snippet")
				a.saveScript([]byte(snippet.GetContent()), result.SnippetTitle, stringutil.StringOrDefault(result.Filename, assistant.RandomScriptFilename()))
			}
		}
	}
}

func (a *appImpl) generateSnippetWithAdditionalPrompt(newPrompt string, history []chat.HistoryEntry, asst assistant.Assistant) {
	log.Debug().Int("prompt_count", len(history)).Msg("Generating additional snippet with assistant")
	a.generateAndHandleScript(newPrompt, history, asst)
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
