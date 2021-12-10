package pathutil

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	fileModeDirectory = os.ModeDir | 0o700
)

// CreatePath returns a suitable location relative to which the file pointed by
// `path` can be written.
func CreatePath(path string) error {
	dir := filepath.Dir(path)
	if Exists(dir) {
		return nil
	}
	if err := os.MkdirAll(dir, fileModeDirectory); err == nil {
		return nil
	}

	return fmt.Errorf("could not create any of the following path: %s", path)
}

// Exists returns true if the specified path exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
