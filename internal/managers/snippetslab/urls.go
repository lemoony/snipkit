package snippetslab

import (
	"bytes"
	"errors"
	"os"
	"path"

	"github.com/phuslu/log"
	"howett.net/plist"

	"github.com/lemoony/snipkit/internal/utils/system"
)

var (
	errNoUserDefinedLibraryPathFound = errors.New("no user defined library path found")
	errNoPreferencesFound            = errors.New("no valid preferences url found")
)

func findLibraryURL(system *system.System, preferencesPath string) snippetsLabLibrary {
	libPath := invalidSnippetsLabLibrary

	if preferencesPath != "" {
		if res, err2 := parsePreferencesForLibraryPath(preferencesPath); err2 == nil {
			libPath = res
		}
	}

	if libPath == invalidSnippetsLabLibrary {
		libPath = snippetsLabLibrary(path.Join(system.UserContainersHome(), defaultPathContaninersLibrary))
	}

	return libPath
}

func findPreferencesPath(system *system.System) (string, error) {
	preferencesURLs, err := getPossiblePreferencesURLs(system)
	if err != nil {
		return "", err
	}
	for _, prefPath := range preferencesURLs {
		fileBytes, err := os.ReadFile(prefPath) //nolint:gosec // potential file inclusion via variable
		if err != nil {
			log.Trace().Msgf("could not open possible preference path %s: %e", prefPath, err)
			continue
		}

		buf := bytes.NewReader(fileBytes)
		decoder := plist.NewDecoder(buf)

		fileMap := make(map[string]interface{})
		if err := decoder.Decode(&fileMap); err != nil {
			log.Trace().Msgf("could not decode possible preference path %s: %e", prefPath, err)
			continue
		}

		if _, ok := fileMap[userDesignatedLibraryPathString]; ok {
			return prefPath, nil
		} else {
			log.Trace().Msgf("invalid preferences file %s: library path not found", prefPath)
			return "", errNoUserDefinedLibraryPathFound
		}
	}

	return "", errNoPreferencesFound
}

func getPossiblePreferencesURLs(system *system.System) ([]string, error) {
	var configDirs []string

	if containerPreferences, err := system.UserContainerPreferences(appID); err != nil {
		return nil, err
	} else {
		configDirs = append(configDirs, path.Join(containerPreferences, preferencesFile))
	}

	return configDirs, nil
}

func parsePreferencesForLibraryPath(preferencesPath string) (snippetsLabLibrary, error) {
	fileBytes, err := os.ReadFile(preferencesPath) //nolint:gosec // potential file inclusion via variable
	if err != nil {
		return invalidSnippetsLabLibrary, err
	}

	buf := bytes.NewReader(fileBytes)
	decoder := plist.NewDecoder(buf)

	fileMap := make(map[string]interface{})
	if err := decoder.Decode(&fileMap); err != nil {
		return invalidSnippetsLabLibrary, err
	}

	if libPath, ok := fileMap[userDesignatedLibraryPathString]; ok {
		return snippetsLabLibrary(libPath.(string)), nil
	}

	return invalidSnippetsLabLibrary, nil
}
