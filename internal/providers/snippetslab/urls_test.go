package snippetslab

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snippet-kit/internal/utils"
)

func Test_getLibraryURL(t *testing.T) {
	system, _ := utils.NewSystem(utils.WithUserHomeDir(testUserHomePath))

	url, err := getLibraryURL(&system)
	assert.NoError(t, err)
	assert.Equal(t, snippetsLabLibrary(testLibraryPath), url)
}

func Test_getPreferencesURL(t *testing.T) {
	system, _ := utils.NewSystem(utils.WithUserHomeDir(testUserHomePath))

	url, err := getPreferencesURL(&system)
	assert.NoError(t, err)
	assert.Equal(t, testPreferencesPath, url.Path)
}
