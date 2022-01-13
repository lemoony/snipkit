package config

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_serializeToYamlWithComment(t *testing.T) {
	var testConfig VersionWrapper
	testConfig.Version = "1.0.0"
	testConfig.Config.Editor = "foo-editor"
	testConfig.Config.Style.Theme = "dracula"
	testConfig.Config.Providers.SnippetsLab.Enabled = true
	testConfig.Config.Providers.SnippetsLab.LibraryPath = "/path/to/lib"
	testConfig.Config.Providers.SnippetsLab.IncludeTags = []string{"snipkit", "othertag"}

	testConfig.Config.Providers.FsLibrary.Enabled = true
	testConfig.Config.Providers.FsLibrary.LibraryPath = []string{"/path/to/file/system/library"}
	testConfig.Config.Providers.FsLibrary.SuffixRegex = []string{".sh"}
	testConfig.Config.Providers.FsLibrary.LazyOpen = true
	testConfig.Config.Providers.FsLibrary.HideTitleInPreview = true

	expectedConfigBytes, err := ioutil.ReadFile(testDataExampleConfig)
	assert.NoError(t, err)

	actualCfgBytes := serializeToYamlWithComment(testConfig)
	assert.Equal(t, string(expectedConfigBytes), string(actualCfgBytes))
}
