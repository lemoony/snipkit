package app

import (
	"os"
	"time"

	"emperror.dev/errors"

	"github.com/lemoony/snipkit/internal/assistant"
	"github.com/lemoony/snipkit/internal/tmpdir"
	"github.com/lemoony/snipkit/internal/ui"
)

func (a *appImpl) CreateSnippetWithAI() {
	if ok, text := a.tui.ShowAiPrompt(); ok {
		stopChan := make(chan bool)

		// Run the spinner in a separate goroutine
		go a.tui.ShowSpinner(text, stopChan)

		asst := assistant.NewBuilder(a.system, a.config.Assistant, a.cache)

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
