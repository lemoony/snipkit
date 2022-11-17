package cmd

import (
	"testing"

	mocks "github.com/lemoony/snipkit/mocks/app"
)

func Test_Browse(t *testing.T) {
	app := mocks.App{}
	app.On("LookupSnippet").Return(nil, nil)

	runExecuteTest(t, []string{"browse"}, withApp(&app))

	app.AssertNumberOfCalls(t, "LookupSnippet", 1)
}
