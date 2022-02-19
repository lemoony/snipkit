package form

import (
	"path/filepath"
	"strings"

	"github.com/phuslu/log"
	"github.com/spf13/afero"

	"github.com/lemoony/snipkit/internal/utils/stringutil"
)

func suggestionsForPath(fs afero.Fs, path string) []string {
	var result []string

	log.Trace().Msgf("Path input: %s", path)

	var isFile bool
	var isPointingToDir bool
	var dirPath string
	var filePath string

	stat, err := fs.Stat(path)
	switch {
	case err == nil && stat.IsDir():
		dirPath = path
		filePath = ""
		isFile = false
		isPointingToDir = true
	case stat != nil:
		dirPath, filePath = filepath.Split(path)
		isFile = true
	default:
		dirPath, filePath = filepath.Split(path)
	}

	log.Trace().Msgf("DirPath: %s, FilePath: %s, isFile: %v, isPointingToDir: %v",
		dirPath, filePath, isFile, isPointingToDir,
	)

	if !isFile {
		if isPointingToDir {
			result = append(result, dirPath)
		}

		if filesInDir, readDirErr := afero.ReadDir(fs, stringutil.StringOrDefault(dirPath, "./")); readDirErr == nil {
			for _, fileInDir := range filesInDir {
				altFilePath := dirPath
				if dirPath != "" && !strings.HasSuffix(altFilePath, "/") {
					altFilePath += "/"
				}
				altFilePath += fileInDir.Name()

				log.Trace().Msgf("alt file path: %s", altFilePath)

				if strings.HasPrefix(altFilePath, path) {
					result = append(result, altFilePath)
				}
			}
		}
	}
	return result
}
