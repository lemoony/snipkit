package utils

import (
	"fmt"
	"os"
	"testing"

	"github.com/adrg/xdg"
	"github.com/stretchr/testify/assert"
)

func Test_System_Default(t *testing.T) {
	s := NewSystem()
	assert.NotNil(t, s)

	assert.Nil(t, s.userDataDir)
	assert.Nil(t, s.userConfigDirs)
	assert.Nil(t, s.userContainersDir)

	assert.NotNil(t, s.UserContainersHome())
	assert.NotEmpty(t, s.UserDataHome())
	assert.NotEmpty(t, s.UserConfigDirs())

	if v, err := s.UserContainerPreferences("test-app"); err != nil {
		assert.NoError(t, err)
		assert.Nil(t, v)
	} else {
		assert.NotNil(t, v)
		assert.Contains(t, v, "test-app")
	}
}

func Test_System_WithOptions(t *testing.T) {
	s := NewSystem(
		WithUserConfigDirs([]string{"/test/config/dir-0", "/test/config/dir-1"}),
		WithUserDataDir("/test/user/data"),
		WithUserContainersDir("/test/container/dir"),
	)
	assert.NotNil(t, s)

	assert.NotNil(t, s.userConfigDirs)

	userConfigDirs := s.UserConfigDirs()
	assert.Equal(t, userConfigDirs[0], "/test/config/dir-0")
	assert.Equal(t, userConfigDirs[1], "/test/config/dir-1")

	assert.Equal(t, s.UserDataHome(), "/test/user/data")
	assert.Equal(t, s.UserContainersHome(), "/test/container/dir")

	if v, err := s.UserContainerPreferences("test-app"); err != nil {
		assert.Nil(t, v)
		assert.NoError(t, err)
	} else {
		assert.Equal(t, v, "/test/container/dir/test-app/Data/Library/Preferences")
	}
}

func Test_getConfigPath(t *testing.T) {
	tests := []struct {
		name              string
		prepare           func()
		system            *System
		expectedConfig    string
		expectedThemesDir string
	}{
		{
			name:              "default",
			system:            NewSystem(),
			prepare:           func() { assert.NoError(t, os.Unsetenv(envSnipkitHome)) },
			expectedConfig:    fmt.Sprintf("%s/snipkit/config.yaml", xdg.ConfigHome),
			expectedThemesDir: fmt.Sprintf("%s/snipkit/themes", xdg.ConfigHome),
		},
		{
			name:              "overwrite config home",
			system:            NewSystem(WithConfigCome("/custom/home")),
			prepare:           func() { assert.NoError(t, os.Unsetenv(envSnipkitHome)) },
			expectedConfig:    "/custom/home/config.yaml",
			expectedThemesDir: "/custom/home/themes",
		},
		{
			name:              "config home via env var",
			system:            NewSystem(),
			prepare:           func() { assert.NoError(t, os.Setenv(envSnipkitHome, "/custom/home/via/env")) },
			expectedConfig:    "/custom/home/via/env/config.yaml",
			expectedThemesDir: "/custom/home/via/env/themes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()
			assert.Equal(t, tt.expectedConfig, tt.system.ConfigPath())
			assert.Equal(t, tt.expectedThemesDir, tt.system.ThemesDir())
		})
	}
}
