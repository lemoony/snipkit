package uimsg

import (
	"testing"

	"github.com/AlecAivazis/survey/v2/core"
	"github.com/stretchr/testify/assert"
)

const (
	testHomePath   = "path/to/tome"
	testCfgPath    = "path/to/config.yml"
	testThemesPath = "path/to/themes"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true

	SetHighlightColor("#ffffff")
}

func Test_ConfigFileCreated(t *testing.T) {
	assert.Contains(t, ConfigFileCreated(testCfgPath), testCfgPath)
}

func Test_ConfigFileDeleted(t *testing.T) {
	assert.Contains(t, ConfigFileDeleted(testCfgPath), testCfgPath)
}

func Test_ConfigNotFound(t *testing.T) {
	assert.Contains(t, ConfigNotFound(testCfgPath), testCfgPath)
}

func Test_SimpleStrings(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func() string
	}{
		{name: "ThemesDeleted", testFunc: ThemesDeleted},
		{name: "ThemesNotDeleted", testFunc: ThemesNotDeleted},
		{name: "ConfigFileCreateConfirm", testFunc: ConfigFileCreateConfirm},
		{name: "ConfigFileRecreateConfirm", testFunc: ConfigFileRecreateConfirm},
		{name: "ThemesDirDeleteConfirm", testFunc: ThemesDirDeleteConfirm},
		{name: "ConfigFileDeleteConfirm", testFunc: ConfigFileDeleteConfirm},
	}

	for _, tt := range tests {
		assert.NotEmpty(t, tt.testFunc())
	}
}

func Test_ConfigFileCreateDescription(t *testing.T) {
	tests := []struct {
		name    string
		homeEnv string
		cfgPath string
	}{
		{
			name:    "home env set",
			homeEnv: "/some/path",
			cfgPath: "/some/path/cfg",
		},
		{
			name:    "home env not set",
			homeEnv: "",
			cfgPath: "/some/other/cfg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := ConfigFileCreateDescription(tt.cfgPath, tt.homeEnv)
			assert.Contains(t, str, tt.cfgPath)
			assert.Contains(t, str, tt.homeEnv)
		})
	}
}

func Test_HomeDirectoryStillExists(t *testing.T) {
	assert.Contains(t, HomeDirectoryStillExists(testHomePath), testHomePath)
}

func Test_ConfigFileRecreateDescription(t *testing.T) {
	assert.Contains(t, ConfigFileRecreateDescription(testHomePath), testHomePath)
}

func Test_ConfigFileDeleteDescription(t *testing.T) {
	assert.Contains(t, ConfigFileDeleteDescription(testHomePath), testHomePath)
}

func Test_ThemesDirDeleteDescription(t *testing.T) {
	assert.Contains(t, ThemesDirDeleteDescription(testThemesPath), testThemesPath)
}

func Test_renderInvalidTemplate(t *testing.T) {
	assert.Panics(t, func() {
		_ = render("{{ if .var }} bla", map[string]interface{}{})
	})
}
