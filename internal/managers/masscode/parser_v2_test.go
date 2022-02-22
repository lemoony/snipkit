package masscode

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/testutil"
)

func Test_parseDBFileV2(t *testing.T) {
	sys := testutil.NewTestSystem()

	snippets := parseDBFileV2(sys, testDataLibraryV2Path)
	assert.Len(t, snippets, 3)

	assert.Equal(t, "176c30e0-2e5d-4be8-a2f2-970eba03901c", snippets[0].GetID())
	assert.Equal(t, "Another", snippets[0].GetTitle())
	assert.Equal(t, model.LanguageText, snippets[0].GetLanguage())
	assert.Equal(t, "echo Hello world", snippets[0].GetContent())

	assert.Equal(t, "Echo something", snippets[1].GetTitle())
	assert.Equal(t, model.LanguageBash, snippets[1].GetLanguage())
	assert.Equal(t, []string{"snipkit"}, snippets[1].GetTags())
	assert.Len(t, snippets[1].GetParameters(), 3)

	assert.Equal(t, "markdown file", snippets[2].GetTitle())
	assert.Equal(t, model.LanguageMarkdown, snippets[2].GetLanguage())
}
