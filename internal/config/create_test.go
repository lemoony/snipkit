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
	testConfig.Config.Providers.SnippetsLab.ExcludeTags = []string{}

	expectedConfigBytes, err := ioutil.ReadFile(testDataExampleConfig)
	assert.NoError(t, err)

	actualCfgBytes, err := serializeToYamlWithComment(testConfig)
	assert.NoError(t, err)
	assert.Equal(t, string(expectedConfigBytes), string(actualCfgBytes))
}
