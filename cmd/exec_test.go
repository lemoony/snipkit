package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snippet-kit/mocks"
)

func Test_Exec(t *testing.T) {
	app := mocks.App{}
	app.On("LookupAndExecuteSnippet").Return(nil)

	err := runMockedTest(t, []string{"exec"}, &app)

	assert.NoError(t, err)
	app.AssertNumberOfCalls(t, "LookupAndExecuteSnippet", 1)
}
