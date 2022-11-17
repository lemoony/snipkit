package cmd

import (
	"testing"

	mocks "github.com/lemoony/snipkit/mocks/app"
)

func Test_Exec(t *testing.T) {
	app := mocks.App{}
	app.On("LookupAndExecuteSnippet", false, false).Return(nil)

	runExecuteTest(t, []string{"exec"}, withApp(&app))

	app.AssertNumberOfCalls(t, "LookupAndExecuteSnippet", 1)
	app.AssertCalled(t, "LookupAndExecuteSnippet", false, false)
}
