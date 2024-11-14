package cmd

import (
	"testing"

	"github.com/stretchr/testify/mock"

	mocks "github.com/lemoony/snipkit/mocks/app"
)

func Test_Assistant_GenerateCmd(t *testing.T) {
	defer resetCommand(execCmd)

	app := mocks.App{}
	app.On("GenerateSnippetWithAssistant", mock.Anything, mock.Anything).Return(nil)

	runExecuteTest(t, []string{"assistant", "generate"}, withApp(&app))

	app.AssertNumberOfCalls(t, "GenerateSnippetWithAssistant", 1)
}

func Test_Assistant_Choose(t *testing.T) {
	defer resetCommand(execCmd)

	app := mocks.App{}
	app.On("EnableAssistant").Return(nil)

	runExecuteTest(t, []string{"assistant", "choose"}, withApp(&app))

	app.AssertNumberOfCalls(t, "EnableAssistant", 1)
}
