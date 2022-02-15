package system

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/adrg/xdg"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/utils/assertutil"
)

func Test_System_Default(t *testing.T) {
	s := NewSystem()
	assert.NotNil(t, s)

	assert.Nil(t, s.userDataDir)
	assert.Nil(t, s.userConfigDirs)
	assert.Nil(t, s.userContainersDir)

	assert.NotNil(t, s.UserContainersHome())
	assert.NotEmpty(t, s.UserHome())
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
		WithUserHome("/test/user"),
		WithUserContainersDir("/test/container/dir"),
	)
	assert.NotNil(t, s)

	assert.NotNil(t, s.userConfigDirs)

	userConfigDirs := s.UserConfigDirs()
	assert.Equal(t, userConfigDirs[0], "/test/config/dir-0")
	assert.Equal(t, userConfigDirs[1], "/test/config/dir-1")

	assert.Equal(t, s.UserHome(), "/test/user")
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

func Test_IsEmpty(t *testing.T) {
	fs := afero.NewMemMapFs()
	path := filepath.Join(t.TempDir(), "foo.txt")

	const dirPerm = 0o755
	assert.NoError(t, fs.Mkdir(filepath.Dir(path), dirPerm))

	system := NewSystem(WithFS(fs))

	assert.True(t, system.IsEmpty(filepath.Dir(path)))

	createTestFile(t, fs, path)
	assert.False(t, system.IsEmpty(filepath.Dir(path)))
}

func Test_DirExists(t *testing.T) {
	fs := afero.NewMemMapFs()
	path := filepath.Join(t.TempDir(), "foo.txt")

	system := NewSystem(WithFS(fs))
	assert.False(t, system.DirExists(filepath.Dir(path)))

	const dirPerm = 0o755
	assert.NoError(t, fs.Mkdir(filepath.Dir(path), dirPerm))

	assert.True(t, system.DirExists(filepath.Dir(path)))
}

func Test_Remove(t *testing.T) {
	fs := afero.NewMemMapFs()
	path := filepath.Join(t.TempDir(), "foo.txt")
	createTestFile(t, fs, path)

	system := NewSystem(WithFS(fs))
	system.Remove(path)

	exists, err := afero.Exists(fs, path)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func Test_RemoveNoPermission(t *testing.T) {
	fs := afero.NewMemMapFs()
	path := filepath.Join(t.TempDir(), "foo.txt")
	createTestFile(t, fs, path)

	system := NewSystem(WithFS(afero.NewReadOnlyFs(fs)))
	_ = assertutil.AssertPanicsWithError(t, ErrFileSystem{}, func() {
		system.Remove(filepath.Dir(path))
	})
}

func Test_RemoveAll(t *testing.T) {
	fs := afero.NewMemMapFs()
	path := filepath.Join(t.TempDir(), "some/path", "foo.txt")
	createTestFile(t, fs, path)

	system := NewSystem(WithFS(fs))
	system.RemoveAll(filepath.Dir(path))

	exists, err := afero.Exists(fs, filepath.Dir(path))
	assert.NoError(t, err)
	assert.False(t, exists)
}

func Test_RemoveAllNoPermission(t *testing.T) {
	fs := afero.NewMemMapFs()

	path := filepath.Join(t.TempDir(), "some/path", "foo.txt")
	createTestFile(t, fs, path)

	system := NewSystem(WithFS(afero.NewReadOnlyFs(fs)))

	_ = assertutil.AssertPanicsWithError(t, ErrFileSystem{}, func() {
		system.RemoveAll(filepath.Dir(path))
	})
}

func createTestFile(t *testing.T, fs afero.Fs, path string) {
	t.Helper()

	const dirPerm = 0o755
	const filePerm = 0o600

	dirPath := filepath.Dir(path)

	assert.NoError(t, os.MkdirAll(dirPath, dirPerm))
	assert.NoError(t, afero.WriteFile(fs, path, []byte("foo"), filePerm))

	exists, err := afero.Exists(fs, path)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func Test_CreatePath(t *testing.T) {
	fs := afero.NewMemMapFs()
	system := NewSystem(WithFS(fs))
	tempDir := t.TempDir()

	testDir := filepath.Join(tempDir, "/foo-dir/testfile")
	testFile := filepath.Join(testDir, "testfile")

	assertutil.AssertExists(t, fs, testDir, false)
	assertutil.AssertExists(t, fs, testFile, false)

	system.CreatePath(testFile)

	assertutil.AssertExists(t, fs, testDir, true)
	assertutil.AssertExists(t, fs, testFile, false)

	system.CreatePath(testFile)
}

func Test_CreatePath_NoPermission(t *testing.T) {
	fs := afero.NewMemMapFs()
	system := NewSystem(WithFS(afero.NewReadOnlyFs(fs)))
	testFile := path.Join(t.TempDir(), "/foo-dir/testfile", "testfile")

	_ = assertutil.AssertPanicsWithError(t, ErrFileSystem{}, func() {
		system.CreatePath(testFile)
	})
}

func Test_WriteFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	system := NewSystem(WithFS(fs))
	filePath := filepath.Join(t.TempDir(), "test.txt")

	assertutil.AssertExists(t, fs, filePath, false)
	system.WriteFile(filePath, []byte("foo"))
	assertutil.AssertExists(t, fs, filePath, true)

	contents, err := afero.ReadFile(fs, filePath)
	assert.NoError(t, err)
	assert.Equal(t, "foo", string(contents))
}

func Test_WriteFileNoPermission(t *testing.T) {
	fs := afero.NewMemMapFs()
	system := NewSystem(WithFS(afero.NewReadOnlyFs(fs)))
	filePath := filepath.Join(t.TempDir(), "test.txt")

	_ = assertutil.AssertPanicsWithError(t, ErrFileSystem{}, func() {
		system.WriteFile(filePath, []byte("foo"))
	})
}

func Test_ReadFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	system := NewSystem(WithFS(fs))
	filePath := filepath.Join(t.TempDir(), "test.txt")

	assert.NoError(t, afero.WriteFile(fs, filePath, []byte("foo"), fileModeConfig))

	assert.Equal(t, "foo", string(system.ReadFile(filePath)))
}

func Test_ReadFile_DoesntExist(t *testing.T) {
	fs := afero.NewMemMapFs()
	filePath := filepath.Join(t.TempDir(), "test.txt")
	system := NewSystem(WithFS(fs))
	_ = assertutil.AssertPanicsWithError(t, ErrFileSystem{}, func() {
		_ = system.ReadFile(filePath)
	})
}

func Test_HomeEnvValue(t *testing.T) {
	_ = os.Setenv(envSnipkitHome, "test")
	assert.Equal(t, "test", NewSystem().HomeEnvValue())
}
