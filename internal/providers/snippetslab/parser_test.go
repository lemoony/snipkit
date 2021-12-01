package snippetslab

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseTags(t *testing.T) {
	tags, err := parseTags(testLibraryPath)
	assert.NoError(t, err)
	assert.Len(t, tags, 1)
	assert.Equal(t, tags["2DA8009E-7BE7-420D-AD57-E7F9BB3ADCBE"], "snipkit")
}

func Test_parseSnippets(t *testing.T) {
	library := snippetsLabLibrary(testLibraryPath)

	snippets, err := parseSnippets(library)
	assert.NoError(t, err)
	assert.Len(t, snippets, 2)
	assert.Equal(t, "Simple echo", snippets[0].Title)
	assert.Equal(t, "Foos script", snippets[1].Title)
}
