package pathutil

import (
	"path"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func Test_Exists(t *testing.T) {
	fs := afero.NewMemMapFs()
	file, err := afero.TempFile(fs, t.TempDir(), "testfile")

	assert.NoError(t, err)
	assert.True(t, Exists(fs, file.Name()))

	assert.NoError(t, file.Close())

	assert.NoError(t, fs.Remove(file.Name()))
	assert.False(t, Exists(fs, file.Name()))
}

func Test_CreatePath(t *testing.T) {
	fs := afero.NewMemMapFs()

	tempDir := t.TempDir()

	testDir := path.Join(tempDir, "/foo-dir/testfile")
	testFile := path.Join(testDir, "testfile")

	assert.False(t, Exists(fs, testDir))
	assert.False(t, Exists(fs, testFile))

	assert.NoError(t, CreatePath(fs, testFile))

	assert.True(t, Exists(fs, testDir))
	assert.False(t, Exists(fs, testFile))

	assert.NoError(t, CreatePath(fs, testFile))
}
