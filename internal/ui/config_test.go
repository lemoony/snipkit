package ui

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/utils/assertutil"
	"github.com/lemoony/snipkit/internal/utils/testutil"
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
		{name: "default", config: DefaultConfig(), previewSchema: "friendly"},
		{name: "default light", config: configWithTheme("default.light"), previewSchema: "friendly"},
		{name: "default dark", config: configWithTheme("default.dark"), previewSchema: "friendly"},
		{name: "simple", config: configWithTheme("simple"), previewSchema: "pastie"},
		{name: "test-custom", config: configWithTheme("test-custom"), previewSchema: "rainbow_dash"},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			theme := tt.config.GetSelectedTheme(system)
			ApplyConfig(tt.config, system)
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
