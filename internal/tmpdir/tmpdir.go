package tmpdir

import (
	"fmt"
	"os"
	"sync"

	"github.com/lemoony/snipkit/internal/utils/system"
)

type TmpDir interface {
	CreateTempFile(contents []byte) (bool, string)
	ClearFiles()
}

type tmpDirImpl struct {
	system  *system.System
	dirPath string
	files   []string
	mutex   sync.Mutex
}

// CreateTempFile creates a temporary file with the provided contents.
// It returns a boolean indicating success, and the path to the created file.
func (t *tmpDirImpl) CreateTempFile(contents []byte) (bool, string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Ensure temporary directory exists
	if err := t.ensureTempDir(); err != nil {
		return false, ""
	}

	// Create a temporary file inside the directory
	tempFile, err := os.CreateTemp(t.dirPath, "tmpfile_*.txt")
	if err != nil {
		return false, ""
	}
	defer func(tempFile *os.File) {
		_ = tempFile.Close()
	}(tempFile)

	// Write contents to the temporary file
	if _, err = tempFile.Write(contents); err != nil {
		return false, ""
	}

	t.files = append(t.files, tempFile.Name())
	return true, tempFile.Name()
}

// ClearFiles removes all files created in the temporary directory.
func (t *tmpDirImpl) ClearFiles() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for _, filePath := range t.files {
		_ = os.Remove(filePath) // Ignore error to continue clearing other files
	}
	t.files = nil
}

// ensureTempDir creates the temporary directory if it does not already exist.
func (t *tmpDirImpl) ensureTempDir() error {
	if t.dirPath == "" {
		dir, err := os.MkdirTemp("", "tmpdir_*")
		if err != nil {
			return fmt.Errorf("failed to create temporary directory: %v", err)
		}
		t.dirPath = dir
	}
	return nil
}

func New(s *system.System) TmpDir {
	return &tmpDirImpl{
		system: s,
		files:  []string{},
	}
}
