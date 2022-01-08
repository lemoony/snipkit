package ui

import (
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func Test_Config_apply(t *testing.T) {
	configWithTheme := func(theme string) Config {
		cfg := DefaultConfig()
		cfg.Theme = theme
		return cfg
	}

	configs := []Config{
		DefaultConfig(),
		configWithTheme("funky"),
		configWithTheme(""),
	}

	for _, cfg := range configs {
		theme := cfg.GetSelectedTheme()
		assert.NotEqual(t, theme.borderColor(), tview.Styles.BorderColor)
		ApplyConfig(cfg)
		assert.Equal(t, theme.borderColor(), tview.Styles.BorderColor)
	}
}
