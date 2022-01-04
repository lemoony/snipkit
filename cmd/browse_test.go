package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snippet-kit/mocks"
)

func Test_Browse(t *testing.T) {
	app := mocks.App{}
	app.On("LookupSnippet").Return(nil, nil)

	err := runMockedTest(t, []string{"browse"}, &app)

	assert.NoError(t, err)
	app.AssertNumberOfCalls(t, "LookupSnippet", 1)
}
