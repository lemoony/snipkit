package masscode

import (
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/idutil"
	"github.com/lemoony/snipkit/internal/utils/testutil"
)

func Test_parseDBFileV1(t *testing.T) {
	sys := testutil.NewTestSystem()

	snippets := parseDBFileV1(sys, testDataMassCodeV1Path)
	assert.Len(t, snippets, 2)

	sort.Slice(snippets, func(i, j int) bool {
		return snippets[i].GetID() < snippets[j].GetID()
	})

	assert.Equal(t, idutil.FormatSnippetID("0fpdzlOnHyoQbSjy", idPrefix), snippets[0].GetID())
	assert.Equal(t, "Echo something", snippets[0].GetTitle())
	assert.Equal(t, model.LanguageBash, snippets[0].GetLanguage())
	assert.Equal(t, []string{"snipkit"}, snippets[0].GetTags())
	assert.Len(t, snippets[0].GetParameters(), 3)
	assert.Equal(t,
		"# some comment\necho \"one\"\n\necho \"two\"\n\necho \"three\"",
		snippets[0].Format([]string{"one", "two", "three"},
			model.SnippetFormatOptions{ParamMode: model.SnippetParamModeReplace, RemoveComments: true}),
	)
}

func Test_parseRawSnippetsV1(t *testing.T) {
	sys := testutil.NewTestSystem()
	snippets := parseRawSnippetsV1(sys, filepath.Join(testDataMassCodeV1Path, v1SnippetsFile))

	assert.Len(t, snippets, 2)

	snippet1 := snippets["a68grSXbYlL5eZkQ"]
	assert.Equal(t, "a68grSXbYlL5eZkQ", snippet1.ID)
	assert.Equal(t, "Simple echo", snippet1.Name)
	assert.Equal(t, "shell", snippet1.Content[0].Language)
	assert.Len(t, snippet1.Content, 1)
	assert.Equal(t, `echo "Hello"`, snippet1.Content[0].Value)
	assert.Empty(t, snippet1.Tags)
	assert.Empty(t, snippet1.TagIDs)

	snippet2 := snippets["0fpdzlOnHyoQbSjy"]
	assert.Equal(t, []string{"HJPPsaOIjc5GaDZe"}, snippet2.Tags)
	assert.Empty(t, snippet2.TagIDs)
}

func Test_parseTagMap(t *testing.T) {
	sys := testutil.NewTestSystem()

	tags := parseRawTagMapV1(sys, filepath.Join(testDataMassCodeV1Path, v1TagsFile))

	assert.Len(t, tags, 1)
	assert.Equal(t, "snipkit", tags["HJPPsaOIjc5GaDZe"])
}
