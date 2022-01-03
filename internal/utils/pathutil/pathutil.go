package pathutil

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

const (
	fileModeDirectory = os.ModeDir | 0o700
)

// CreatePath returns a suitable location relative to which the file pointed by
// `path` can be written.
func CreatePath(fs afero.Fs, path string) error {
	dir := filepath.Dir(path)
	if ok, err := afero.DirExists(fs, path); err != nil {
		panic(err)
	} else if ok {
		return nil
	}

	if err := fs.MkdirAll(dir, fileModeDirectory); err == nil {
		return nil
	}

	return fmt.Errorf("could not create any of the following path: %s", path)
}

// Exists returns true if the specified path exists.
func Exists(fs afero.Fs, path string) bool {
	ok, err := afero.Exists(fs, path)
	if err != nil {
		panic(err)
	}
	return ok
}
