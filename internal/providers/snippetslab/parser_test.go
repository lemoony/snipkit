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
	assert.Equal(t, model.LanguageBash, snippets[0].GetLanguage())
}
