package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_System_Default(t *testing.T) {
	s, err := NewSystem()
	assert.NoError(t, err)
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
	s, err := NewSystem(
		WithUserConfigDirs([]string{"/test/config/dir-0", "/test/config/dir-1"}),
		WithUserDataDir("/test/user/data"),
		WithUserContainersDir("/test/container/dir"),
	)
	assert.NoError(t, err)
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
