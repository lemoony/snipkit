package pictarinesnip

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/stringutil"
	"github.com/lemoony/snipkit/internal/utils/testutil"
)

func Test_parseLibrary(t *testing.T) {
	system := testutil.NewTestSystem()
	snippets := parseLibrary(testDataDefaultLibraryPath, system, &stringutil.StringSet{})

	assert.Len(t, snippets, 2)

	assert.Equal(t, "Echo something", snippets[0].GetTitle())
	assert.Regexp(t, "^# some comment.*", snippets[0].GetContent())
	assert.Equal(t, model.LanguageBash, snippets[0].GetLanguage())

	assert.Equal(t, "Another snippet", snippets[1].GetTitle())
	assert.Equal(t, "echo \"Hello\"", snippets[1].GetContent())
	assert.Equal(t, model.LanguageUnknown, snippets[1].GetLanguage())
}
