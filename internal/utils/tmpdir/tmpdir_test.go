package tmpdir

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/utils/testutil"
)

func TestCreateTempFile(t *testing.T) {
	// Initialize a new tmpDir instance
	sys := testutil.NewTestSystem()
	tmpDir := New(sys)

	// Test content to write into the temporary file
	content := []byte("test content")
	success, path := tmpDir.CreateTempFile(content)

	// Assert file creation was successful
	assert.True(t, success)
	assert.NotEmpty(t, path)

	// Read file contents to confirm it matches
	data, err := afero.ReadFile(sys.Fs, path)
	assert.NoError(t, err)
	assert.Equal(t, content, data)

	// Cleanup
	tmpDir.ClearFiles()
}

func TestClearFiles(t *testing.T) {
	sys := testutil.NewTestSystem()
	tmpDir := New(sys)

	// Create several temporary files
	for i := 0; i < 3; i++ {
		content := []byte("test content")
		success, _ := tmpDir.CreateTempFile(content)
		assert.True(t, success)
	}

	// Clear the created files
	tmpDir.ClearFiles()

	x := tmpDir.(*tmpDirImpl)

	assert.False(t, x.system.DirExists(x.dirPath))
}
