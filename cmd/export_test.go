package cmd

import (
	"testing"

	"github.com/stretchr/testify/mock"

	mocks "github.com/lemoony/snipkit/mocks/app"
)

func Test_Export(t *testing.T) {
	app := mocks.App{}
	app.On("ExportSnippets", mock.AnythingOfType("[]app.ExportField"), mock.AnythingOfType("app.ExportFormat")).
		Return("{}")

	runExecuteTest(t, []string{"export", "-f=id,title", "--output", "json-pretty"}, withApp(&app))

	app.AssertNumberOfCalls(t, "ExportSnippets", 1)
}
