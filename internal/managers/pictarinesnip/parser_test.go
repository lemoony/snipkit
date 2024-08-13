package pictarinesnip

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/idutil"
	"github.com/lemoony/snipkit/internal/utils/stringutil"
	"github.com/lemoony/snipkit/internal/utils/testutil"
)

func Test_parseLibrary(t *testing.T) {
	system := testutil.NewTestSystem()
	snippets := parseLibrary(testDataDefaultLibraryPath, system, &stringutil.StringSet{})

	assert.Len(t, snippets, 2)

	snippet1 := snippets[0]
	assert.Equal(t, idutil.FormatSnippetID("88235A31-F0AD-4206-96DE-19E0EDEE79B2", idPrefix), snippet1.GetID())
	assert.Equal(t, "Echo something", snippet1.GetTitle())
	assert.Regexp(t, "^# some comment.*", snippet1.GetContent())
	assert.Equal(t, model.LanguageBash, snippet1.GetLanguage())
	assert.Equal(t, []string{"snipkit"}, snippet1.GetTags())
	assert.Len(t, snippet1.GetParameters(), 3)
	assert.NotEqual(t, snippet1.GetContent(), snippet1.Format([]string{"one", "two", "three"}, model.SnippetFormatOptions{}))

	snippet2 := snippets[1]
	assert.Equal(t, idutil.FormatSnippetID("B3473DF8-6ED6-4589-BFFC-C75F73B1B522", idPrefix), snippet2.GetID())
	assert.Equal(t, "Another snippet", snippet2.GetTitle())
	assert.Equal(t, "echo \"Hello\"", snippet2.GetContent())
	assert.Equal(t, model.LanguageUnknown, snippet2.GetLanguage())
	assert.Equal(t, []string{}, snippet2.GetTags())
	assert.Empty(t, snippet2.GetParameters())
	assert.Equal(t, snippet2.GetContent(), snippet2.Format([]string{}, model.SnippetFormatOptions{}))
}
