package ui

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snippet-kit/internal/utils/testutil"
)

func Test_defaultThemeAvailable(t *testing.T) {
	themeNames := []string{
		defaultThemeName,
		"funky",
	}
	for _, themeName := range themeNames {
		t.Run(themeName, func(t *testing.T) {
			assert.NotNil(t, embeddedTheme(themeName))
		})
	}
}

func Test_themeNotFound(t *testing.T) {
	err := testutil.AssertPanicsWithError(t, ErrInvalidTheme, func() {
		embeddedTheme("foo-theme")
	})

	assert.Contains(t, err.Error(), "theme not found: foo-theme")
}

func Test_embeddedThemes(t *testing.T) {
	fileInfos, err := ioutil.ReadDir("../../themes")
	assert.NoError(t, err)

	for _, fileInfo := range fileInfos {
		if filepath.Ext(fileInfo.Name()) == ".go" {
			continue
		}

		themeName := strings.TrimSuffix(fileInfo.Name(), ".yaml")
		assert.NotNil(t, embeddedTheme(themeName))
	}
}
