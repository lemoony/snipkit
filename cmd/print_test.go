package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snippet-kit/mocks"
)

func Test_Print(t *testing.T) {
	app := mocks.App{}
	app.On("LookupAndCreatePrintableSnippet").
		Return("snippet-printed", nil)

	err := runMockedTest(t, []string{"print"}, &app)

	assert.NoError(t, err)
	app.AssertNumberOfCalls(t, "LookupAndCreatePrintableSnippet", 1)
}
