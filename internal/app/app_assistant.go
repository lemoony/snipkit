package app

import (
	"fmt"
	"os"
	"strings"
	"time"

	"emperror.dev/errors"

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

func (a *appImpl) GenerateSnippetWithAssistant(demoScriptPath string, demoQueryDuration time.Duration) {
	asst := a.assistantProviderFunc(a.config.Assistant)
	if valid, msg := asst.Initialize(); !valid {
		a.tui.PrintAndExit(msg, -1)
	}

	if ok, text := a.tui.ShowAssistantPrompt([]string{}); ok {
		prompts := []string{text}

		spinnerStop := a.startSpinner()
		script, filename := a.getScriptWithAssistant(asst, demoScriptPath, demoQueryDuration, text)
		close(spinnerStop)

		a.handleGeneratedScript(script, filename, prompts, asst)
	}
}

func (a *appImpl) getScriptWithAssistant(asst assistant.Assistant, demoScriptPath string, demoQueryDuration time.Duration, prompt string) (string, string) {
	if demoScriptPath != "" {
		demoScript := a.system.ReadFile(demoScriptPath)
		time.Sleep(demoQueryDuration)
		return string(demoScript), "demo.sh"
	}
	return asst.Query(prompt)
}

func (a *appImpl) handleGeneratedScript(script, filename string, prompts []string, asst assistant.Assistant) {
	tmpDirSvc := tmpdir.New(a.system)
	defer tmpDirSvc.ClearFiles()

	if fileOk, filePath := tmpDirSvc.CreateTempFile([]byte(script)); fileOk {
		a.tui.OpenEditor(filePath, a.config.Editor)
		//nolint:gosec // ignore potential file inclusion via variable
		if updatedContents, err := os.ReadFile(filePath); err != nil {
			panic(errors.Wrapf(err, "failed to read temporary file"))
		} else {
			snippet := assistant.PrepareSnippet(updatedContents)
			var parameterValues []string
			paramOk := true
			if parameters := snippet.GetParameters(); len(parameters) > 0 {
				parameterValues, paramOk = a.tui.ShowParameterForm(snippet.GetParameters(), nil, ui.OkButtonExecute)
			}
			if paramOk {
				a.executeAndHandleSnippet(snippet, parameterValues, prompts, asst, filename)
			}
		}
	}
}

func (a *appImpl) executeAndHandleSnippet(snippet model.Snippet, parameterValues []string, prompts []string, asst assistant.Assistant, filename string) {
	if executed, capturedResult := a.executeSnippet(false, false, snippet, parameterValues); executed {
		if result := a.tui.ShowAssistantWizard(wizard.Config{ProposedFilename: filename}); result.SelectedOption == wizard.OptionTryAgain {
			if ok2, prompt2 := a.tui.ShowAssistantPrompt(prompts); ok2 {
				prompts = append(prompts, prompt2)
				newPrompt := fmt.Sprintf("The result of the command was: %s\n%s\n\n%s", capturedResult.stdout, capturedResult.stderr, prompt2)
				a.generateSnippetWithAdditionalPrompt(newPrompt, prompts, asst)
			}
		} else if result.SelectedOption == wizard.OptionSaveExit {
			a.saveScript([]byte(snippet.GetContent()), snippet.GetTitle(), stringutil.StringOrDefault(result.Filename, assistant.RandomScriptFilename()))
		}
	}
}

func (a *appImpl) generateSnippetWithAdditionalPrompt(newPrompt string, prompts []string, asst assistant.Assistant) {
	spinnerStop := a.startSpinner()
	script, filename := asst.Query(newPrompt)
	close(spinnerStop)

	a.handleGeneratedScript(script, filename, prompts, asst)
}

func (a *appImpl) startSpinner() chan bool {
	stopChan := make(chan bool)
	go a.tui.ShowSpinner("Please wait, generating script...", stopChan)
	return stopChan
}

func (a *appImpl) saveScript(contents []byte, title, filename string) {
	if manager, ok := a.getSaveAssistantSnippetHelper(); ok {
		manager.SaveAssistantSnippet(title, filename, contents)
	}
}

func (a *appImpl) getSaveAssistantSnippetHelper() (managers.Manager, bool) {
	if manager, ok := sliceutil.FindElement(a.managers, func(manager managers.Manager) bool {
		return manager.Key() == fslibrary.Key
	}); ok {
		return manager, true
	} else {
		panic("File system library not configured as manager. Try run `snipkit manager add`")
	}
}

func (a *appImpl) EnableAssistant() {
	asst := a.assistantProviderFunc(a.config.Assistant)

	assistantDescriptions := asst.AssistantDescriptions(a.config.Assistant)
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
		cfg := asst.AutoConfig(assistantDescription.Key)
		configBytes := config.SerializeToYamlWithComment(cfg)
		configStr := strings.TrimSpace(string(configBytes))
		confirmed := a.tui.Confirmation(uimsg.ManagerConfigAddConfirm(configStr))
		if confirmed {
			a.configService.UpdateAssistantConfig(cfg)
		}
		a.tui.Print(uimsg.AssistantUpdateConfigResult(confirmed, a.configService.ConfigFilePath()))
	}
}
