package ui

import (
	"path/filepath"
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snippet-kit/internal/utils/assertutil"
	"github.com/lemoony/snippet-kit/internal/utils/testutil"
)

func Test_Config_apply(t *testing.T) {
	configWithTheme := func(theme string) Config {
		cfg := DefaultConfig()
		cfg.Theme = theme
		return cfg
	}

	system := testutil.NewTestSystem()

	bytes := system.ReadFile("testdata/test-custom.yaml")
	path := filepath.Join(system.ThemesDir(), "test-custom.yaml")
	system.CreatePath(path)
	system.WriteFile(path, bytes)

	testdata := []struct {
		name          string
		config        Config
		previewSchema string
	}{
		{name: "default theme", config: DefaultConfig(), previewSchema: "friendly"},
		{name: "example theme", config: configWithTheme("example"), previewSchema: "solarized-light"},
		{name: "test-custom", config: configWithTheme("test-custom"), previewSchema: "monokai"},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			theme := tt.config.GetSelectedTheme(system)
			assert.NotEqual(t, theme.borderColor(), tview.Styles.BorderColor)
			ApplyConfig(tt.config, system)
			assert.Equal(t, theme.borderColor(), tview.Styles.BorderColor)
			assert.Equal(t, tt.previewSchema, theme.PreviewColorSchemeName)
		})
	}
}

func Test_GetUnknownTheme(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Theme = "foo-theme"

	system := testutil.NewTestSystem()

	err := assertutil.AssertPanicsWithError(t, ErrInvalidTheme, func() {
		cfg.GetSelectedTheme(system)
	})

	assert.Contains(t, err.Error(), "theme not found: "+cfg.Theme)
}
