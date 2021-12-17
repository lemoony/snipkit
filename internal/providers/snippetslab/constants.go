package snippetslab

import (
	"net/url"
	"path"
)

type snippetsLabLibrary string

func (t *snippetsLabLibrary) basePath() (string, error) {
	libURL, err := url.Parse(string(*t))
	if err != nil {
		return "", err
	}
	return libURL.Path, nil
}

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

func (t *snippetsLabLibrary) validate() (bool, error) {
	if _, err := parseTags(*t); err != nil {
		return false, err
	}
	return true, nil
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

	appID           = "com.renfei.SnippetsLab"
	preferencesFile = appID + ".plist"

	databaseSubPath = "Database"
	tagsSubPath     = databaseSubPath + "/tags.data"
	snippetsSubPath = databaseSubPath + "/Snippets"

	defaultPathContaninersLibrary   = appID + "/Data/Library/Application Support/" + appID + "/main.snippetslablibrary"
	userDesignatedLibraryPathString = "User DesignatedLibraryPathString"

	invalidSnippetsLabLibrary = snippetsLabLibrary("")
)
