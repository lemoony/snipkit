package snippetslab

import (
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

	for _, s := range snippets {
		then.AssertThat(t,
			s.GetTitle(),
			is.AnyOf(is.EqualTo("Simple echo"), is.EqualTo("Foos script")),
		)
		then.AssertThat(t,
			s.GetContent(),
			is.AnyOf(is.MatchForPattern("^# some comment.*"), is.MatchForPattern("echo \"Foo!\"")),
		)
	}
	assert.Equal(t, "84A08C4A-B2BE-4964-A521-180550BDA7B3", snippets[0].GetID())
	assert.Empty(t, snippets[0].GetTags())
	assert.Equal(t, model.LanguageBash, snippets[0].GetLanguage())
	assert.Len(t, snippets[0].GetParameters(), 2)
	assert.NotEqual(t, snippets[0].Format([]string{"one", "two"}), snippets[0].GetContent())

	assert.Equal(t, []string{"2DA8009E-7BE7-420D-AD57-E7F9BB3ADCBE"}, snippets[1].GetTags())
}
