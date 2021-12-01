package snippetslab

import (
	"bytes"
	"io/ioutil"
	"net/url"
	"path"

	"howett.net/plist"

	"github.com/lemoony/snippet-kit/internal/utils"
)

func getLibraryURL(system *utils.System) (snippetsLabLibrary, error) {
	if preferencesURL, err := getPreferencesURL(system); err != nil {
		return invalidSnippetsLabLibrary, err
	} else if libURL, err := parsePreferencesForLibraryURL(system, preferencesURL); err != nil {
		return invalidSnippetsLabLibrary, err
	} else {
		return snippetsLabLibrary(libURL.Path), nil
	}
}

func getPreferencesURL(system *utils.System) (*url.URL, error) {
	homeDir, err := system.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return url.Parse(path.Join(homeDir, defaultSettingsPath))
}

func parsePreferencesForLibraryURL(system *utils.System, preferencesURL *url.URL) (*url.URL, error) {
	currentUserHomeDir, err := system.UserHomeDir()
	if err != nil {
		return nil, err
	}

	fileBytes, err := ioutil.ReadFile(preferencesURL.Path)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewReader(fileBytes)
	decoder := plist.NewDecoder(buf)

	fileMap := make(map[string]interface{})
	if err := decoder.Decode(&fileMap); err != nil {
		return nil, err
	}

	if path, ok := fileMap[userDesignatedLibraryPathString]; ok {
		return url.Parse(path.(string))
	}

	if homeDir, err := url.Parse(currentUserHomeDir); err != nil {
		return nil, err
	} else {
		return url.Parse(path.Join(homeDir.Path, defaultLibraryPath))
	}
}
