package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	mocks "github.com/lemoony/snippet-kit/mocks/app"
)

func Test_Exec(t *testing.T) {
	app := mocks.App{}
	app.On("LookupAndExecuteSnippet").Return(nil)

	err := runMockedTest(t, []string{"exec"}, withApp(&app))

	assert.NoError(t, err)
	app.AssertNumberOfCalls(t, "LookupAndExecuteSnippet", 1)
}
