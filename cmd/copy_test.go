package cmd

import (
	"testing"

	mocks "github.com/lemoony/snipkit/mocks/app"
)

func Test_Copy(t *testing.T) {
	app := mocks.App{}
	app.On("LookupAndCreatePrintableSnippet").
		Return(true, "snippet-printed")

	runExecuteTest(t, []string{"copy"}, withApp(&app))

	app.AssertNumberOfCalls(t, "LookupAndCreatePrintableSnippet", 1)

	assertClipboardContent(t, "snippet-printed")
}
