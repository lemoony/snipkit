package snippetslab

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/utils/system"
)

func Test_getLibraryURLFromPreferencesFile(t *testing.T) {
	system := system.NewSystem(system.WithUserContainersDir(testdataContainersPath))

	url := findLibraryURL(system, testDataPreferencesWithUserDefinedLibraryPath)
	assert.Equal(t, snippetsLabLibrary("file://"+testDataDefaultLibraryPath), url)
}

func Test_getLibraryURLDefault(t *testing.T) {
	system := system.NewSystem(system.WithUserContainersDir(testdataContainersPath))

	url := findLibraryURL(system, testDataPreferencesPath)
	assert.Equal(t, snippetsLabLibrary(testDataDefaultLibraryPath), url)
}

func Test_parsePreferencesForLibraryPath(t *testing.T) {
	// fileMap := make(map[string]interface{})
	// fileMap[userDesignatedLibraryPathString] = "file://" + testDataDefaultLibraryPath
	// f, err := os.Create(testDataPreferencesWithUserDefinedLibraryPath)
	// encoder := plist.NewEncoderForFormat(f, plist.BinaryFormat)
	// encoder.Encode(fileMap)
	libURL, err := parsePreferencesForLibraryPath(testDataPreferencesWithUserDefinedLibraryPath)
	assert.NoError(t, err)
	assert.Equal(t, "file://"+testDataDefaultLibraryPath, string(libURL))
}

func Test_parsePreferencesForLibraryPath_Notound(t *testing.T) {
	library, err := parsePreferencesForLibraryPath("path/does/not/exist/by/purpose")
	assert.Equal(t, invalidSnippetsLabLibrary, library)
	assert.Error(t, err)
}

func Test_getPreferencesURL(t *testing.T) {
	system := system.NewSystem(system.WithUserContainersDir(testdataContainersPath))

	urls, err := getPossiblePreferencesURLs(system)
	assert.NoError(t, err)
	assert.Len(t, urls, 1)
	assert.Equal(t, testDataPreferencesPath, urls[0])
}

func Test_findPreferencesPath_NoUserDefinedPath(t *testing.T) {
	// fileMap := make(map[string]interface{})
	// fileMap[userDesignatedLibraryPathString] = "file://" + testDataDefaultLibraryPath
	// f, _ := os.Create(testDataPreferencesPath)
	// encoder := plist.NewEncoderForFormat(f, plist.BinaryFormat)
	// encoder.Encode(fileMap)
	system := system.NewSystem(system.WithUserContainersDir(testdataContainersPath))
	path, err := findPreferencesPath(system)
	assert.Error(t, err)
	assert.ErrorIs(t, errNoUserDefinedLibraryPathFound, err)
	assert.Equal(t, "", path)
}
