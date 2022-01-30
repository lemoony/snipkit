package app

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snipkit/internal/config/configtest"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/assertutil"
	"github.com/lemoony/snipkit/internal/utils/testutil/mockutil"
	uiMocks "github.com/lemoony/snipkit/mocks/ui"
)

func Test_Lookup_ErrNoSnippetsAvailable(t *testing.T) {
	var snippets []model.Snippet

	tui := uiMocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	tui.On("ShowLookup", mock.Anything).Return(1)

	app := NewApp(
		WithTUI(&tui), WithConfig(configtest.NewTestConfig().Config), withManagerSnippets(snippets),
	)

	_ = assertutil.AssertPanicsWithError(t, ErrNoSnippetsAvailable, func() {
		app.LookupSnippet()
	})
}
