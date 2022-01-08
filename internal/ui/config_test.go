package ui

import (
	"path/filepath"
	"testing"

	"github.com/rivo/tview"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snippet-kit/internal/utils"
	"github.com/lemoony/snippet-kit/internal/utils/testutil"
)

func Test_Config_apply(t *testing.T) {
	configWithTheme := func(theme string) Config {
		cfg := DefaultConfig()
		cfg.Theme = theme
		return cfg
	}

	system := utils.NewTestSystem()

	if bytes, err := afero.ReadFile(system.Fs, "testdata/test-custom.yaml"); err != nil {
		assert.NoError(t, err)
	} else {
		assert.NoError(
			t,
			afero.WriteFile(system.Fs, filepath.Join(system.ThemesDir(), "test-custom.yaml"), bytes, 0o600),
		)
	}

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

	system := utils.NewTestSystem()

	err := testutil.AssertPanicsWithError(t, ErrInvalidTheme, func() {
		cfg.GetSelectedTheme(system)
	})

	assert.Contains(t, err.Error(), "theme not found: "+cfg.Theme)
}
