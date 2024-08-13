package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	appx "github.com/lemoony/snipkit/internal/app"
	mocks "github.com/lemoony/snipkit/mocks/app"
)

func Test_Export(t *testing.T) {
	app := mocks.App{}
	app.On("ExportSnippets", mock.AnythingOfType("[]app.ExportField"), mock.AnythingOfType("app.ExportFormat")).
		Return("{}")

	runExecuteTest(t, []string{"export", "-f=id,title", "--output", "json-pretty"}, withApp(&app))

	app.AssertNumberOfCalls(t, "ExportSnippets", 1)

	exportFields := app.Calls[0].Arguments.Get(0).([]appx.ExportField)
	assert.Len(t, exportFields, 2)
	assert.Equal(t, appx.ExportFieldID, exportFields[0])
	assert.Equal(t, appx.ExportFieldTitle, exportFields[1])

	assert.Equal(t, appx.ExportFormatPrettyJSON, app.Calls[0].Arguments.Get(1).(appx.ExportFormat))
}
