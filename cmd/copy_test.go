package cmd

import (
	"os"
	"runtime"
	"testing"

	"github.com/atotto/clipboard"
	"github.com/stretchr/testify/assert"

	mocks "github.com/lemoony/snipkit/mocks/app"
)

func Test_Copy(t *testing.T) {
	app := mocks.App{}
	app.On("LookupAndCreatePrintableSnippet").
		Return("snippet-printed", true)

	runExecuteTest(t, []string{"copy"}, withApp(&app))

	app.AssertNumberOfCalls(t, "LookupAndCreatePrintableSnippet", 1)

	// Workaround: Clipboard doesn't work on the CI for linux.
	if runtime.GOOS != "linux" || os.Getenv("CI") != "true" {
		if content, err := clipboard.ReadAll(); err != nil {
			assert.NoError(t, err)
		} else {
			assert.Equal(t, "snippet-printed", content)
		}
	}
}
