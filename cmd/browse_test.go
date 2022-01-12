package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	mocks "github.com/lemoony/snipkit/mocks/app"
)

func Test_Browse(t *testing.T) {
	app := mocks.App{}
	app.On("LookupSnippet").Return(nil, nil)

	err := runMockedTest(t, []string{"browse"}, withApp(&app))

	assert.NoError(t, err)
	app.AssertNumberOfCalls(t, "LookupSnippet", 1)
}
