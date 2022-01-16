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
	testConfig.Config.Manager.SnippetsLab.Enabled = true
	testConfig.Config.Manager.SnippetsLab.LibraryPath = "/path/to/lib"
	testConfig.Config.Manager.SnippetsLab.IncludeTags = []string{"snipkit", "othertag"}

	testConfig.Config.Manager.FsLibrary.Enabled = true
	testConfig.Config.Manager.FsLibrary.LibraryPath = []string{"/path/to/file/system/library"}
	testConfig.Config.Manager.FsLibrary.SuffixRegex = []string{".sh"}
	testConfig.Config.Manager.FsLibrary.LazyOpen = true
	testConfig.Config.Manager.FsLibrary.HideTitleInPreview = true

	expectedConfigBytes, err := ioutil.ReadFile(testDataExampleConfig)
	assert.NoError(t, err)

	actualCfgBytes := serializeToYamlWithComment(testConfig)
	assert.Equal(t, string(expectedConfigBytes), string(actualCfgBytes))
}
