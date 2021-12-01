package snippetslab

import (
	"net/url"
	"path"
)

type snippetsLabLibrary string

func (t *snippetsLabLibrary) tagsFilePath() (string, error) {
	if libURL, err := url.Parse(string(*t)); err != nil {
		return "", err
	} else {
		return path.Join(libURL.Path, tagsSubPath), nil
	}
}

func (t *snippetsLabLibrary) snippetsFilePath() (string, error) {
	if libURL, err := url.Parse(string(*t)); err != nil {
		return "", err
	} else {
		return path.Join(libURL.Path, snippetsSubPath), nil
	}
}

const (
	SnippetTitle        = "com.renfei.SnippetsLab.Key.SnippetTitle"
	SnippetParts        = "com.renfei.SnippetsLab.Key.SnippetParts"
	SnippetPartContent  = "com.renfei.SnippetsLab.Key.SnippetPartContent"
	SnippetPartLanguage = "com.renfei.SnippetsLab.Key.SnippetPartLanguage"
	SnippetUUID         = "com.renfei.SnippetsLab.Key.SnippetUUID"
	SnippetTagUUIDs     = "com.renfei.SnippetsLab.Key.SnippetTagUUIDs"

	SnippetTagsTagUUID  = "com.renfei.SnippetsLab.Key.TagUUID"
	SnippetTagsTagTitle = "com.renfei.SnippetsLab.Key.TagTitle"

	databaseSubPath = "Database"
	tagsSubPath     = databaseSubPath + "/tags.data"
	snippetsSubPath = databaseSubPath + "/Snippets"

	defaultLibraryPath              = "Library/Containers/com.renfei.SnippetsLab/Data/Library/Application Support/com.renfei.SnippetsLab/main.snippetslablibrary"
	defaultSettingsPath             = "Library/Containers/com.renfei.SnippetsLab/Data/Library/Preferences/com.renfei.SnippetsLab.plist"
	userDesignatedLibraryPathString = "User DesignatedLibraryPathString"

	invalidSnippetsLabLibrary = snippetsLabLibrary("")
)
