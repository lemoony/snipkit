package snippetslab

import (
	"sort"
	"testing"

	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/model"
)

func Test_parseTags(t *testing.T) {
	tags, err := parseTags(testDataDefaultLibraryPath)
	assert.NoError(t, err)
	assert.Len(t, tags, 1)
	assert.Equal(t, tags["2DA8009E-7BE7-420D-AD57-E7F9BB3ADCBE"], "snipkit")
}

func Test_parseSnippets(t *testing.T) {
	library := snippetsLabLibrary(testDataDefaultLibraryPath)

	snippets, err := parseSnippets(library)
	assert.NoError(t, err)
	assert.Len(t, snippets, 2)

	sort.Slice(snippets, func(i, j int) bool {
		return snippets[i].GetID() < snippets[j].GetID()
	})

	snippet1 := snippets[0]
	assert.Equal(t, "84A08C4A-B2BE-4964-A521-180550BDA7B3", snippet1.GetID())
	assert.Empty(t, snippet1.GetTags())
	assert.Equal(t, model.LanguageBash, snippet1.GetLanguage())
	assert.Len(t, snippet1.GetParameters(), 2)
	assert.NotEqual(t, snippet1.Format([]string{"one", "two"}), snippet1.GetContent())
	then.AssertThat(t, snippet1.GetContent(), is.MatchForPattern("^# some comment.*"))
	then.AssertThat(t, snippet1.GetTitle(), is.AnyOf(is.EqualTo("Simple echo")))

	snippet2 := snippets[1]
	assert.Equal(t, "B3EDC3BE-6FE1-489E-9EB8-C400D4CF1B54", snippet2.GetID())
	assert.Equal(t, []string{"2DA8009E-7BE7-420D-AD57-E7F9BB3ADCBE"}, snippet2.GetTags())
	assert.Equal(t, model.LanguageBash, snippet2.GetLanguage())
	assert.Empty(t, snippet2.GetParameters())
	assert.NotEqual(t, snippet2.Format([]string{}), snippet1.GetContent())
	then.AssertThat(t, snippet2.GetContent(), is.MatchForPattern("echo \"Foo!\""))
	then.AssertThat(t, snippet2.GetTitle(), is.AnyOf(is.EqualTo("Foos script")))
}
