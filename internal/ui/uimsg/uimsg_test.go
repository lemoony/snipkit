package uimsg

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/ui/style"
	"github.com/lemoony/snipkit/internal/utils/testutil"
)

const (
	testHomePath   = "path/to/tome"
	testCfgPath    = "path/to/config.yml"
	testThemesPath = "path/to/themes"
)

var testStyle = style.NoopStyle

func Test_ConfigFileCreated(t *testing.T) {
	assert.Contains(t, render(ConfigFileCreateResult(true, testCfgPath, false)), testCfgPath)
}

func Test_ConfigNotFound(t *testing.T) {
	assert.Contains(t, render(ConfigNotFound(testCfgPath)), testCfgPath)
}

func Test_ConfigNeedsMigration(t *testing.T) {
	assert.Contains(t, render(ConfigNeedsMigration("1.0", "2.0")), "migrate the config file")
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
			assert.Contains(t, testutil.StripANSI(c.Header(testStyle, 0)), tt.cfgPath)
			assert.Contains(t, testutil.StripANSI(c.Header(testStyle, 0)), tt.homeEnv)
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
	assert.Contains(t, testutil.StripANSI(c.Header(testStyle, 0)), testCfgPath)
}

func Test_ConfigFileDeleteResult(t *testing.T) {
	assert.Contains(t, render(ConfigFileDeleteResult(true, testCfgPath)), "Configuration file deleted")
}

func Test_ConfigFileMigrationConfirm(t *testing.T) {
	configStr := `config: test`
	c := ConfigFileMigrationConfirm(configStr)
	assert.Contains(t, testutil.StripANSI(c.Header(testStyle, 0)), configStr)
}

func Test_ConfigFileMigrationResult(t *testing.T) {
	assert.Contains(
		t,
		render(ConfigFileMigrationResult(false, testCfgPath)),
		"The configuration file was not updated.",
	)
}

func Test_ExecConfirm(t *testing.T) {
	c := ExecConfirm("test-title", "print hello")
	result := testutil.StripANSI(c.Header(testStyle, 0))
	assert.Contains(t, result, "Snippet: test-title")
	assert.Contains(t, result, "print hello")
}

func Test_ExecPrint(t *testing.T) {
	c := ExecPrint("title", "print hello")
	assert.Contains(t, render(c), "Snippet: title")
}

func Test_HomeDirectoryStillExists(t *testing.T) {
	assert.Contains(t, render(HomeDirectoryStillExists(testHomePath)), testHomePath)
}

func Test_ThemesDeleteConfirm(t *testing.T) {
	c := ThemesDeleteConfirm(testThemesPath)
	assert.NotEmpty(t, c.Prompt)
	assert.Contains(t, testutil.StripANSI(c.Header(testStyle, 0)), testThemesPath)
}

func Test_ThemesDeleteResult(t *testing.T) {
	assert.Contains(t, render(ThemesDeleteResult(true, testThemesPath)), testThemesPath)
}

func Test_ManagerConfigAddConfirm(t *testing.T) {
	c := ManagerConfigAddConfirm("yaml")
	assert.NotEmpty(t, c.Prompt)
	assert.Contains(t, c.Header(testStyle, 0), "yaml")
}

func Test_ManagerAddConfigResult(t *testing.T) {
	assert.Contains(t, render(ManagerAddConfigResult(true, testCfgPath)), testCfgPath)
}

func Test_ManagerOauthDeviceFlow(t *testing.T) {
	assert.Contains(t, render(ManagerOauthDeviceFlow("github.com", "1234-5678")), "1234-5678")
}

func Test_renderInvalidTemplate(t *testing.T) {
	assert.Panics(t, func() {
		_ = renderWithStyle("{{ if .var }} bla", testStyle, map[string]interface{}{})
	})
}

func render(p Printable) string {
	return testutil.StripANSI(p.RenderWith(testStyle))
}
