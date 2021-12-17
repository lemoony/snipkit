package config

import (
	"io/ioutil"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const testDataExampleConfig = "testdata/example-config.yaml"

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

func Test_loadConfigFromYAL(t *testing.T) {
	v := viper.New()
	v.SetConfigFile(testDataExampleConfig)

	config, err := LoadConfig(v)

	assert.NoError(t, err)
	assert.NotNil(t, config)

	assert.Equal(t, "foo-editor", config.Editor)
	assert.Equal(t, "dracula", config.Style.Theme)
	assert.True(t, config.Providers.SnippetsLab.Enabled)
	assert.Equal(t, "/path/to/lib", config.Providers.SnippetsLab.LibraryPath)
	assert.Len(t, config.Providers.SnippetsLab.IncludeTags, 2)
	assert.Len(t, config.Providers.SnippetsLab.ExcludeTags, 0)
}
