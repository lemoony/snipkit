package snippetslab

import (
	"testing"

	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"github.com/stretchr/testify/assert"
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

	then.AssertThat(t,
		snippets[0].Title,
		is.AnyOf(is.EqualTo("Simple echo"), is.EqualTo("Foos script")),
	)
}
