package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	mocks "github.com/lemoony/snipkit/mocks/app"
)

func Test_Exec(t *testing.T) {
	app := mocks.App{}
	app.On("LookupAndExecuteSnippet", false).Return(nil)

	err := runMockedTest(t, []string{"exec"}, withApp(&app))

	assert.NoError(t, err)
	app.AssertNumberOfCalls(t, "LookupAndExecuteSnippet", 1)
	app.AssertCalled(t, "LookupAndExecuteSnippet", false)
}
