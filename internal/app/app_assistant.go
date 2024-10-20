package app

import (
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
	"github.com/lemoony/snipkit/internal/ui/picker"
	"github.com/lemoony/snipkit/internal/ui/uimsg"
	"github.com/lemoony/snipkit/internal/utils/sliceutil"
	"github.com/lemoony/snipkit/internal/utils/stringutil"
	"github.com/lemoony/snipkit/internal/utils/tmpdir"
)

const (
	saveYes = "yes"
	saveNo  = "no"
)

func (a *appImpl) GenerateSnippetWithAssistant() {
	asst := a.assistantProviderFunc(a.config.Assistant)

	if ok, text := a.tui.ShowPrompt("What do you want the script to do?"); ok {
		stopChan := make(chan bool)

		// Run the spinner in a separate goroutine
		go a.tui.ShowSpinner(text, stopChan)

		script, filename := asst.Query(text)

		// Send stop signal to stop the spinner
		stopChan <- true

		//nolint:mnd // Wait briefly to ensure spinner quits cleanly
		time.Sleep(100 * time.Millisecond)

		tmpDirSvc := tmpdir.New(a.system)
		defer tmpDirSvc.ClearFiles()

		if fileOk, filePath := tmpDirSvc.CreateTempFile([]byte(script)); fileOk {
			a.tui.OpenEditor(filePath, a.config.Editor)
			//nolint:gosec // ignore potential file inclusion via variable
			if updatedContents, err := os.ReadFile(filePath); err != nil {
				panic(errors.Wrapf(err, "failed to read temporary file"))
			} else {
				snippet := assistant.PrepareSnippet(updatedContents)
				parameters := snippet.GetParameters()

				if a.config.Assistant.SaveMode == assistant.SaveModeFsLibrary {
					parameters = append(parameters, saveFsLibParameter())
				}

				if parameterValues, paramOk := a.tui.ShowParameterForm(parameters, nil, ui.OkButtonExecute); paramOk {
					if shouldSaveScript(a.config.Assistant.SaveMode, parameterValues) {
						defer a.saveScript(updatedContents, stringutil.StringOrDefault(filename, assistant.RandomScriptFilename()))
					}
					a.executeSnippet(false, false, snippet, parameterValues)
				}
			}
		}
	}
}

func (a *appImpl) saveScript(contents []byte, filename string) {
	if manager, ok := a.getSaveAssistantSnippetHelper(); ok {
		manager.SaveAssistantSnippet(filename, contents)
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

func shouldSaveScript(saveMode assistant.SaveMode, parameterValues []string) bool {
	return saveMode == assistant.SaveModeFsLibrary && parameterValues[len(parameterValues)-1] == saveYes
}

func saveFsLibParameter() model.Parameter {
	return model.Parameter{
		Key:          "SAVE_FS_LIBRARY",
		Name:         "Save in file system library",
		Description:  "Should be saved to file system library",
		Type:         model.ParameterTypeValue,
		Values:       []string{saveYes, saveNo},
		DefaultValue: saveNo,
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
		asst.AutoConfig(assistantDescription.Key, a.system)
		cfg := asst.AutoConfig(assistantDescription.Key, a.system)
		configBytes := config.SerializeToYamlWithComment(cfg)
		configStr := strings.TrimSpace(string(configBytes))
		confirmed := a.tui.Confirmation(uimsg.ManagerConfigAddConfirm(configStr))
		if confirmed {
			a.configService.UpdateAssistantConfig(cfg)
		}
		a.tui.Print(uimsg.AssistantUpdateConfigResult(confirmed, a.configService.ConfigFilePath()))
	}
}
