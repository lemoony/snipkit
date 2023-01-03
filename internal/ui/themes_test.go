package ui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_defaultThemeAvailable(t *testing.T) {
	themeNames := []string{
		defaultThemeName,
	}
	for _, themeName := range themeNames {
		t.Run(themeName, func(t *testing.T) {
			theme, ok := embeddedTheme(themeName)
			assert.True(t, ok)
			assert.NotNil(t, theme)
		})
	}
}

func Test_themeNotFound(t *testing.T) {
	theme, ok := embeddedTheme("foo-theme")
	assert.Nil(t, theme)
	assert.False(t, ok)
}

func Test_embeddedThemes(t *testing.T) {
	fileInfos, err := os.ReadDir("../../themes")
	assert.NoError(t, err)

	for _, fileInfo := range fileInfos {
		if filepath.Ext(fileInfo.Name()) == ".go" {
			continue
		}

		themeName := strings.TrimSuffix(fileInfo.Name(), ".yaml")
		theme, ok := embeddedTheme(themeName)
		assert.True(t, ok)
		assert.NotNil(t, theme)
	}
}
