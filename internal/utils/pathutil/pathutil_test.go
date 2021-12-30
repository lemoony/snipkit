package pathutil

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Exists(t *testing.T) {
	file, err := ioutil.TempFile(t.TempDir(), "testfile")
	assert.NoError(t, err)
	assert.True(t, Exists(file.Name()))

	assert.NoError(t, os.Remove(file.Name()))
	assert.False(t, Exists(file.Name()))
}

func Test_CreatePath(t *testing.T) {
	tempDir := t.TempDir()

	testDir := path.Join(tempDir, "/foo-dir/testfile")
	testFile := path.Join(testDir, "testfile")

	assert.False(t, Exists(testDir))
	assert.False(t, Exists(testFile))

	assert.NoError(t, CreatePath(testFile))

	assert.True(t, Exists(testDir))
	assert.False(t, Exists(testFile))

	assert.NoError(t, CreatePath(testFile))
}
