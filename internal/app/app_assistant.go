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

	if ok, text := a.tui.ShowAssistantPrompt([]string{}); ok {
		prompts := []string{text}
		spinnerStop := a.startSpinner()
		script := asst.Query(text)
		close(spinnerStop)
		a.handleGeneratedScript(script, prompts, asst)
	}
}

func (a *appImpl) handleGeneratedScript(parsed assistant.ParsedScript, prompts []string, asst assistant.Assistant) {
	tmpDirSvc := tmpdir.New(a.system)
	defer tmpDirSvc.ClearFiles()

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
				a.executeAndHandleSnippet(snippet, parameterValues, prompts, asst, parsed)
			}
		} else {
			a.executeAndHandleSnippet(snippet, nil, prompts, asst, parsed)
		}
	}
}

func (a *appImpl) executeAndHandleSnippet(snippet model.Snippet, parameterValues []string, prompts []string, asst assistant.Assistant, script assistant.ParsedScript) {
	executed, capturedResult := a.executeSnippet(len(parameterValues) == 0, false, snippet, parameterValues)
	if executed {
		wizardOk, result := a.tui.ShowAssistantWizard(wizard.Config{
			ShowSaveOption:      a.config.Assistant.SaveMode != assistant.SaveModeNever,
			ProposedFilename:    script.Filename,
			ProposedSnippetName: script.Title,
		})
		if wizardOk {
			switch result.SelectedOption {
			case wizard.OptionTryAgain:
				log.Debug().Msg("User requested to try again with assistant")
				if ok2, prompt2 := a.tui.ShowAssistantPrompt(prompts); ok2 {
					prompts = append(prompts, prompt2)
					newPrompt := fmt.Sprintf("The result of the command was: %s\n%s\n\n%s", capturedResult.stdout, capturedResult.stderr, prompt2)
					a.generateSnippetWithAdditionalPrompt(newPrompt, prompts, asst)
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

func (a *appImpl) generateSnippetWithAdditionalPrompt(newPrompt string, prompts []string, asst assistant.Assistant) {
	log.Debug().Int("prompt_count", len(prompts)+1).Msg("Generating additional snippet with assistant")
	spinnerStop := a.startSpinner()
	parsed := asst.Query(newPrompt)
	close(spinnerStop)
	a.handleGeneratedScript(parsed, prompts, asst)
}

func (a *appImpl) startSpinner() chan bool {
	stopChan := make(chan bool)
	go a.tui.ShowSpinner("Please wait, generating script...", "SnipKit Assistant", stopChan)
	return stopChan
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

		cfg := assistantInstance.AutoConfig(assistantDescription.Key)
		configBytes := config.SerializeToYamlWithComment(cfg)
		configStr := strings.TrimSpace(string(configBytes))
		confirmed := a.tui.Confirmation(uimsg.ManagerConfigAddConfirm(configStr))
		if confirmed {
			a.configService.UpdateAssistantConfig(cfg)
			log.Debug().Str("provider", assistantDescription.Name).Msg("Assistant configuration updated")
		}
		a.tui.Print(uimsg.AssistantUpdateConfigResult(confirmed, a.configService.ConfigFilePath()))
	}
}
