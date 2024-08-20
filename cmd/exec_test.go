package cmd

import (
	"testing"

	"github.com/lemoony/snipkit/internal/model"
	mocks "github.com/lemoony/snipkit/mocks/app"
)

func Test_Exec(t *testing.T) {
	defer execCmd.SetContext(nil) //nolint:staticcheck // allow nil as context in order to reset

	app := mocks.App{}
	app.On("LookupAndExecuteSnippet", false, false).Return(nil)

	runExecuteTest(t, []string{"exec"}, withApp(&app))

	app.AssertNumberOfCalls(t, "LookupAndExecuteSnippet", 1)
	app.AssertCalled(t, "LookupAndExecuteSnippet", false, false)
}

func Test_Exec_WithFlags(t *testing.T) {
	defer execCmd.SetContext(nil) //nolint:staticcheck // allow nil as context in order to reset

	app := mocks.App{}
	app.On(
		"FindScriptAndExecuteWithParameters",
		"foo",
		[]model.ParameterValue{{Key: "KEY1", Value: "VALUE1"}, {Key: "KEY2", Value: "VALUE2"}},
		false,
		false,
	).Return(nil)

	runExecuteTest(t, []string{"exec", "--id", "foo", "--param", "KEY1=VALUE1", "--param=KEY2=VALUE2"}, withApp(&app))

	app.AssertNumberOfCalls(t, "FindScriptAndExecuteWithParameters", 1)
}
