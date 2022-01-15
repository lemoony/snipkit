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
			c := ConfirmConfigCreation(tt.cfgPath, tt.homeEnv)
			assert.Contains(t, c.Header(), tt.cfgPath)
			assert.Contains(t, c.Header(), tt.homeEnv)
		})
	}
}

func Test_HomeDirectoryStillExists(t *testing.T) {
	assert.Contains(t, HomeDirectoryStillExists(testHomePath), testHomePath)
}

func Test_ConfigFileRecreateDescription(t *testing.T) {
	c := ConfirmConfigRecreate(testCfgPath)
	assert.NotEmpty(t, c.Prompt)
	assert.Contains(t, c.Header(), testCfgPath)
}

func Test_ConfigFilDelete(t *testing.T) {
	c := ConfirmConfigDelete(testCfgPath)
	assert.NotEmpty(t, c.Prompt)
	assert.Contains(t, c.Header(), testCfgPath)
}

func Test_ThemesDirDeleteDescription(t *testing.T) {
	c := ConfirmThemesDelete(testThemesPath)
	assert.NotEmpty(t, c.Prompt)
	assert.Contains(t, c.Header(), testThemesPath)
}

func Test_renderInvalidTemplate(t *testing.T) {
	assert.Panics(t, func() {
		_ = render("{{ if .var }} bla", map[string]interface{}{})
	})
}
