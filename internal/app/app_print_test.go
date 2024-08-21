package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snipkit/internal/config/configtest"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/assertutil"
	"github.com/lemoony/snipkit/internal/utils/testutil"
	"github.com/lemoony/snipkit/internal/utils/testutil/mockutil"
	uiMocks "github.com/lemoony/snipkit/mocks/ui"
)

var expectedPrintOutput = `# ${VAR1} Name: First Output
# ${VAR1} Description: What to print on the tui first
VAR1="foo-value"
echo "${VAR1}`

func Test_LookupAndCreatePrintableSnippet(t *testing.T) {
	snippets := []model.Snippet{
		testutil.TestSnippet{ID: "uuid1", Title: "title-1", Language: model.LanguageYAML, Tags: []string{}, Content: "content-1"},
		testutil.TestSnippet{ID: "uuid2", Title: "title-2", Language: model.LanguageBash, Tags: []string{}, Content: testSnippetContent},
	}

	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	tui.On("ShowLookup", mock.Anything, mock.Anything).Return(1)
	tui.On("ShowParameterForm", mock.Anything, mock.Anything, mock.Anything).Return([]string{"foo-value"}, true)

	app := NewApp(
		WithTUI(&tui), WithConfig(configtest.NewTestConfig().Config), withManagerSnippets(snippets),
	)

	ok, s := app.LookupAndCreatePrintableSnippet()
	assert.True(t, ok)
	assert.Equal(t, expectedPrintOutput, s)
}

func Test_LookupAndCreatePrintableSnippet_NoneSelected(t *testing.T) {
	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	tui.On("ShowLookup", mock.Anything, mock.Anything).Return(-1)
	app := NewApp(
		WithTUI(&tui), WithConfig(configtest.NewTestConfig().Config), withManagerSnippets([]model.Snippet{
			testutil.DummySnippet,
		}),
	)

	ok, _ := app.LookupAndCreatePrintableSnippet()
	assert.False(t, ok)
}

func Test_FindSnippetAndPrint(t *testing.T) {
	snippets := []model.Snippet{
		testutil.TestSnippet{ID: "uuid1", Title: "title-1", Language: model.LanguageYAML, Tags: []string{}, Content: "content-1"},
		testutil.TestSnippet{ID: "uuid2", Title: "title-2", Language: model.LanguageBash, Tags: []string{}, Content: testSnippetContent},
	}

	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()

	app := NewApp(
		WithTUI(&tui), WithConfig(configtest.NewTestConfig().Config), withManagerSnippets(snippets),
	)

	ok, s := app.FindSnippetAndPrint("uuid2", []model.ParameterValue{{Key: "VAR1", Value: "foo-value"}})
	assert.True(t, ok)
	assert.Equal(t, expectedPrintOutput, s)
}

func Test_FindSnippetAndPrint_NotFound(t *testing.T) {
	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()

	app := NewApp(
		WithTUI(&tui), WithConfig(configtest.NewTestConfig().Config), withManagerSnippets([]model.Snippet{
			testutil.DummySnippet,
		}),
	)

	_ = assertutil.AssertPanicsWithError(t, ErrSnippetIDNotFound, func() {
		app.FindSnippetAndPrint("random-id", []model.ParameterValue{})
	})
}

func Test_FindSnippetAndPrint_MissingParameters(t *testing.T) {
	snippets := []model.Snippet{
		testutil.TestSnippet{ID: "uuid1", Title: "title-1", Language: model.LanguageYAML, Tags: []string{}, Content: "content-1"},
		testutil.TestSnippet{ID: "uuid2", Title: "title-2", Language: model.LanguageBash, Tags: []string{}, Content: testSnippetContent},
	}

	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	tui.On("ShowParameterForm", mock.Anything, mock.Anything, mock.Anything).Return([]string{"foo-value"}, true)

	app := NewApp(
		WithTUI(&tui), WithConfig(configtest.NewTestConfig().Config), withManagerSnippets(snippets),
	)

	ok, s := app.FindSnippetAndPrint("uuid2", []model.ParameterValue{})
	assert.True(t, ok)
	assert.Equal(t, expectedPrintOutput, s)
}

func Test_LookupAndPrintSnippetArgs(t *testing.T) {
	snippets := []model.Snippet{
		testutil.TestSnippet{ID: "uuid-x", Title: "title-2", Language: model.LanguageBash, Tags: []string{}, Content: testSnippetContent},
	}

	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	tui.On("ShowLookup", mock.Anything, mock.Anything).Return(0)
	tui.On("ShowParameterForm", mock.Anything, mock.Anything, mock.Anything).Return([]string{"foo-value"}, true)

	app := NewApp(
		WithTUI(&tui), WithConfig(configtest.NewTestConfig().Config), withManagerSnippets(snippets),
	)

	ok, id, parameterValues := app.LookupSnippetArgs()
	assert.True(t, ok)
	assert.Equal(t, id, "uuid-x")
	assert.Equal(t, parameterValues, []model.ParameterValue{{Key: "VAR1", Value: "foo-value"}})
}

func Test_LookupAndPrintSnippetArgs_NoneSelected(t *testing.T) {
	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	tui.On("ShowLookup", mock.Anything, mock.Anything).Return(-1)

	app := NewApp(
		WithTUI(&tui), WithConfig(configtest.NewTestConfig().Config), withManagerSnippets([]model.Snippet{
			testutil.DummySnippet,
		}),
	)

	ok, _, _ := app.LookupSnippetArgs()
	assert.False(t, ok)
}
