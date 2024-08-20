package cmd

import (
	"testing"

	"github.com/lemoony/snipkit/internal/model"
	mocks "github.com/lemoony/snipkit/mocks/app"
)

func Test_Print(t *testing.T) {
	defer resetCommand(printCmd)

	app := mocks.App{}
	app.On("LookupAndCreatePrintableSnippet").
		Return("snippet-printed", true)

	runExecuteTest(t, []string{"print", "--copy"}, withApp(&app))

	app.AssertNumberOfCalls(t, "LookupAndCreatePrintableSnippet", 1)
	assertClipboardContent(t, "snippet-printed")
}

func Test_Print_WithCmdFlags(t *testing.T) {
	defer resetCommand(printCmd)

	app := mocks.App{}
	app.On("FindSnippetAndPrint", "foo", []model.ParameterValue{{Key: "KEY1", Value: "VALUE1"}, {Key: "KEY2", Value: "VALUE2"}}).
		Return("snippet-printed", true)

	runExecuteTest(t, []string{"print", "--copy", "--id", "foo", "--param", "KEY1=VALUE1", "--param=KEY2=VALUE2"}, withApp(&app))

	app.AssertNumberOfCalls(t, "FindSnippetAndPrint", 1)
	assertClipboardContent(t, "snippet-printed")
}
