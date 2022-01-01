package snippetslab

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snippet-kit/internal/utils"
)

func Test_AutoDiscoveryConfig_NoSnippetLab(t *testing.T) {
	system := utils.NewSystem(utils.WithUserContainersDir("path/does/not/exist"))
	config := AutoDiscoveryConfig(&system)
	assert.False(t, config.Enabled)
}

func Test_AutoDiscoveryConfig_Available(t *testing.T) {
	system := utils.NewSystem(utils.WithUserContainersDir(testdataContainersPath))
	config := AutoDiscoveryConfig(&system)
	assert.True(t, config.Enabled)
	assert.Equal(t, testDataDefaultLibraryPath, config.LibraryPath)
}
