package uimsg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testPath = "path/to/config.yml"

func Test_ConfigFileCreated(t *testing.T) {
	assert.Contains(t, ConfigFileCreated(testPath), testPath)
}

func Test_ConfigFileDeleted(t *testing.T) {
	assert.Contains(t, ConfigFileDeleted(testPath), testPath)
}

func Test_ConfigNotFound(t *testing.T) {
	assert.Contains(t, ConfigNotFound(testPath), testPath)
}

func Test_SimpleStrings(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func() string
	}{
		{name: "ConfigNotDeleted", testFunc: ConfigNotDeleted},
		{name: "ThemesDeleted", testFunc: ThemesDeleted},
		{name: "ThemesNotDeleted", testFunc: ThemesNotDeleted},
		{name: "ConfigFileCreateConfirm", testFunc: func() string { return ConfigFileCreateConfirm("path/to") }},
		{name: "ConfigFileRecreateConfirm", testFunc: ConfigFileRecreateConfirm},
		{name: "ConfigFileDeleteConfirm", testFunc: func() string { return ConfigFileDeleteConfirm("path/to") }},
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
	path := "path/to/home"
	assert.Contains(t, HomeDirectoryStillExists(path), path)
}

func Test_ConfigFileRecreateDescription(t *testing.T) {
	path := "path/to/home"
	assert.Contains(t, ConfigFileRecreateDescription(path), path)
}

func Test_ThemesDirDeleteConfirm(t *testing.T) {
	path := "path/to/themes"
	assert.Contains(t, ThemesDirDeleteConfirm(path), path)
}

func Test_renderInvalidTemplate(t *testing.T) {
	assert.Panics(t, func() {
		_ = render("{{ if .var }} bla", map[string]interface{}{})
	})
}
