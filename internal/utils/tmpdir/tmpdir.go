package tmpdir

import (
	"fmt"
	"sync"

	"github.com/spf13/afero"

	"github.com/lemoony/snipkit/internal/utils/system"
)

type TmpDir interface {
	CreateTempFile(contents []byte) (bool, string)
	ClearFiles()
}

type tmpDirImpl struct {
	system  *system.System
	dirPath string
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
	tempFile, err := afero.TempFile(t.system.Fs, t.dirPath, "tmpfile_*.txt")
	if err != nil {
		return false, ""
	}

	defer func(tempFile afero.File) {
		_ = tempFile.Close()
	}(tempFile)

	// Write contents to the temporary file
	if _, err = tempFile.Write(contents); err != nil {
		return false, ""
	}

	return true, tempFile.Name()
}

// ClearFiles removes all files created in the temporary directory.
func (t *tmpDirImpl) ClearFiles() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.system.RemoveAll(t.dirPath)
}

// ensureTempDir creates the temporary directory if it does not already exist.
func (t *tmpDirImpl) ensureTempDir() error {
	if t.dirPath == "" {
		dir, err := afero.TempDir(t.system.Fs, "", "snipkit_tmpdir")
		if err != nil {
			return fmt.Errorf("failed to create temporary directory: %v", err)
		}
		t.dirPath = dir
	}
	return nil
}

func New(s *system.System) TmpDir {
	return &tmpDirImpl{system: s}
}
