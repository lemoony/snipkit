package uimsg

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/ui/style"
)

const (
	testHomePath   = "path/to/tome"
	testCfgPath    = "path/to/config.yml"
	testThemesPath = "path/to/themes"
)

var testStyle = style.NoopStyle

func Test_ConfigFileCreated(t *testing.T) {
	assert.Contains(t, ConfigFileCreateResult(true, testCfgPath, false).RenderWith(testStyle), testCfgPath)
}

func Test_ConfigNotFound(t *testing.T) {
	assert.Contains(t, ConfigNotFound(testCfgPath).RenderWith(testStyle), testCfgPath)
}

func Test_ConfigFileCreateConfirm(t *testing.T) {
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
			c := ConfigFileCreateConfirm(tt.cfgPath, tt.homeEnv, true)
			assert.Contains(t, c.Header(testStyle, 0), tt.cfgPath)
			assert.Contains(t, c.Header(testStyle, 0), tt.homeEnv)
		})
	}
}

func Test_ConfigFileCreateResult(t *testing.T) {
	tests := []struct {
		deleted  bool
		recreate bool
	}{
		{deleted: true, recreate: true},
		{deleted: false, recreate: true},
		{deleted: true, recreate: false},
		{deleted: false, recreate: false},
	}
	for _, tt := range tests {
		c := ConfigFileCreateResult(tt.deleted, testCfgPath, tt.recreate)
		// TODO: assert more
		assert.NotEmpty(t, c)
	}
}

func Test_ConfigFileDeleteConfirm(t *testing.T) {
	c := ConfigFileDeleteConfirm(testCfgPath)
	assert.Contains(t, c.Header(testStyle, 0), testCfgPath)
}

func Test_ConfigFileDeleteResult(t *testing.T) {
	assert.Contains(t, ConfigFileDeleteResult(true, testCfgPath).RenderWith(testStyle), "Configuration file deleted")
}

func Test_HomeDirectoryStillExists(t *testing.T) {
	assert.Contains(t, HomeDirectoryStillExists(testHomePath).RenderWith(testStyle), testHomePath)
}

func Test_ThemesDeleteConfirm(t *testing.T) {
	c := ThemesDeleteConfirm(testThemesPath)
	assert.NotEmpty(t, c.Prompt)
	assert.Contains(t, c.Header(testStyle, 0), testThemesPath)
}

func Test_ThemesDeleteResult(t *testing.T) {
	assert.Contains(t, ThemesDeleteResult(true, testThemesPath).RenderWith(testStyle), testThemesPath)
}

func Test_ManagerConfigAddConfirm(t *testing.T) {
	c := ManagerConfigAddConfirm("yaml")
	assert.NotEmpty(t, c.Prompt)
	assert.Contains(t, c.Header(testStyle, 0), "yaml")
}

func Test_ManagerAddConfigResult(t *testing.T) {
	assert.Contains(t, ManagerAddConfigResult(true, testCfgPath).RenderWith(testStyle), testCfgPath)
}

func Test_renderInvalidTemplate(t *testing.T) {
	assert.Panics(t, func() {
		_ = renderWithStyle("{{ if .var }} bla", testStyle, map[string]interface{}{})
	})
}
