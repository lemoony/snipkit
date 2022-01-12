package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	mocks "github.com/lemoony/snipkit/mocks/app"
)

func Test_Print(t *testing.T) {
	app := mocks.App{}
	app.On("LookupAndCreatePrintableSnippet").
		Return("snippet-printed", true)

	err := runMockedTest(t, []string{"print"}, withApp(&app))

	assert.NoError(t, err)
	app.AssertNumberOfCalls(t, "LookupAndCreatePrintableSnippet", 1)
}
