package app

import (
	"os"
	"strings"
	"time"

	"emperror.dev/errors"

	"github.com/lemoony/snipkit/internal/assistant"
	"github.com/lemoony/snipkit/internal/config"
	"github.com/lemoony/snipkit/internal/tmpdir"
	"github.com/lemoony/snipkit/internal/ui"
	"github.com/lemoony/snipkit/internal/ui/picker"
	"github.com/lemoony/snipkit/internal/ui/uimsg"
)

func (a *appImpl) CreateSnippetWithAI() {
	asst := assistant.NewBuilder(a.system, a.config.Assistant, a.cache)

	if ok, text := a.tui.ShowPrompt("What do you want the script to do?"); ok {
		stopChan := make(chan bool)

		// Run the spinner in a separate goroutine
		go a.tui.ShowSpinner(text, stopChan)

		response := asst.Query(text)

		// Send stop signal to stop the spinner
		stopChan <- true

		//nolint:mnd // Wait briefly to ensure spinner quits cleanly
		time.Sleep(100 * time.Millisecond)

		tmpDirSvc := tmpdir.New(a.system)
		defer tmpDirSvc.ClearFiles()

		if fileOk, filePath := tmpDirSvc.CreateTempFile([]byte(response)); fileOk {
			a.tui.OpenEditor(filePath, a.config.Editor)
			//nolint:gosec // ignore potential file inclusion via variable
			if updatedContents, err := os.ReadFile(filePath); err != nil {
				panic(errors.Wrapf(err, "failed to read temporary file"))
			} else {
				snippet := assistant.PrepareSnippet(string(updatedContents))
				parameters := snippet.GetParameters()
				if parameterValues, paramOk := a.tui.ShowParameterForm(parameters, nil, ui.OkButtonExecute); paramOk {
					a.executeSnippet(false, false, snippet, parameterValues)
				}
			}
		}
	}
}

func (a *appImpl) EnableAssistant() {
	assistant := assistant.NewBuilder(a.system, a.config.Assistant, a.cache)

	assistantDescriptions := assistant.AssistantDescriptions(a.config.Assistant)

	listItems := make([]picker.Item, len(assistantDescriptions))
	var selected *picker.Item
	for i := range assistantDescriptions {
		listItems[i] = picker.NewItem(assistantDescriptions[i].Name, assistantDescriptions[i].Description)
		if assistantDescriptions[i].Enabled {
			selected = &listItems[i]
		}
	}

	if selectedIndex, ok := a.tui.ShowPicker("Which assistant AI do you want to enable?", listItems, selected); ok {
		assistantDescription := assistantDescriptions[selectedIndex]
		assistant.AutoConfig(assistantDescription.Key, a.system)
		cfg := assistant.AutoConfig(assistantDescription.Key, a.system)
		configBytes := config.SerializeToYamlWithComment(cfg)
		configStr := strings.TrimSpace(string(configBytes))
		confirmed := a.tui.Confirmation(uimsg.ManagerConfigAddConfirm(configStr))
		if confirmed {
			a.configService.UpdateAssistantConfig(cfg)
		}
		a.tui.Print(uimsg.ManagerAddConfigResult(confirmed, a.configService.ConfigFilePath()))
	}
}
