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
	oldConfigStr := `config: old`
	newConfigStr := `config: new`
	c := ConfigFileMigrationConfirm(oldConfigStr, newConfigStr)
	result := testutil.StripANSI(c.Header(testStyle, 120))
	assert.NotEmpty(t, result)
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
	oldConfig := "old: yaml"
	newConfig := "new: yaml"
	c := ManagerConfigAddConfirm(oldConfig, newConfig)
	assert.NotEmpty(t, c.Prompt)
	result := c.Header(testStyle, 120)
	assert.NotEmpty(t, result)
}

func Test_ManagerAddConfigResult(t *testing.T) {
	assert.Contains(t, render(ManagerAddConfigResult(true, testCfgPath)), testCfgPath)
}

func Test_ManagerOauthDeviceFlow(t *testing.T) {
	assert.Contains(t, render(ManagerOauthDeviceFlow("github.com", "1234-5678")), "1234-5678")
}

func Test_MAssistantUpdateConfigResult(t *testing.T) {
	assert.Contains(t, render(AssistantUpdateConfigResult(true, testCfgPath)), testCfgPath)
}

func Test_renderInvalidTemplate(t *testing.T) {
	assert.Panics(t, func() {
		_ = renderWithStyle("{{ if .var }} bla", testStyle, map[string]interface{}{})
	})
}

func render(p Printable) string {
	return testutil.StripANSI(p.RenderWith(testStyle))
}

func Test_SideBySideDiff_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		oldYaml string
		newYaml string
		verify  func(t *testing.T, result string)
	}{
		{
			name:    "minimal old config",
			oldYaml: "version: v1.1.0\nconfig: {}",
			newYaml: "version: v1.1.0\nconfig:\n  editor: vim",
			verify: func(t *testing.T, result string) {
				t.Helper()
				// Should render as "new config only" table
				assert.Contains(t, result, "NEW CONFIGURATION")
				assert.Contains(t, result, "editor: vim")
			},
		},
		{
			name:    "empty old config",
			oldYaml: "",
			newYaml: "version: v1.1.0\nconfig:\n  editor: nvim",
			verify: func(t *testing.T, result string) {
				t.Helper()
				// Should render as "new config only" table
				assert.Contains(t, result, "NEW CONFIGURATION")
			},
		},
		{
			name:    "normal diff",
			oldYaml: "version: v1.1.0\nconfig:\n  editor: vim\n  fuzzy: true\n  theme: default",
			newYaml: "version: v1.1.0\nconfig:\n  editor: nvim\n  fuzzy: true\n  theme: dark",
			verify: func(t *testing.T, result string) {
				t.Helper()
				// Should render side-by-side
				assert.Contains(t, result, "BEFORE")
				assert.Contains(t, result, "AFTER")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use the template function directly
			funcs := templateFuncs(testStyle)
			sideBySideDiff := funcs["SideBySideDiff"].(func(...interface{}) string)

			result := sideBySideDiff(tt.oldYaml, tt.newYaml)
			tt.verify(t, testutil.StripANSI(result))
		})
	}
}
