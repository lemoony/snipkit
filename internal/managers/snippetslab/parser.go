package snippetslab

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"howett.net/plist"

	"github.com/lemoony/snipkit/internal/model"
)

const (
	nullEntry = "$null"
)

var languageMapping = map[string]model.Language{
	"YamlLexer":     model.LanguageYAML,
	"BashLexer":     model.LanguageBash,
	"MarkdownLexer": model.LanguageMarkdown,
	"TOMLLexer":     model.LanguageTOML,
	"TextLexer":     model.LanguageText,
}

//nolint:forcetypeassert // since we will catch any panic error and checking each statement explicitly is too much work
func parseTags(library snippetsLabLibrary) (map[string]string, error) {
	path, err := library.tagsFilePath()
	if err != nil {
		return nil, err
	}

	fileMap, err := readPblistFile(path)
	if err != nil {
		return nil, err
	}

	result := map[string]string{}

	objects := fileMap["$objects"].([]interface{})
	var indexOfNull int
	for i, v := range objects {
		if v == nullEntry {
			indexOfNull = i + 1
			break
		}
	}

	keyMapping := objects[indexOfNull].(map[string]interface{})
	tagIndices := keyMapping["NS.objects"].([]interface{})

	for _, v := range tagIndices {
		index := uint64(v.(plist.UID))
		tagMap := objects[index].(map[string]interface{})

		tagMapData := tagMap["NS.data"].([]uint8)

		tagBuffer := bytes.NewReader(tagMapData)
		tagDecoder := plist.NewDecoder(tagBuffer)

		tagFields := make(map[string]interface{})
		if err := tagDecoder.Decode(&tagFields); err != nil {
			return nil, err
		}

		tagObjects := tagFields["$objects"].([]interface{})

		var indexOfNull int
		for i, v := range tagObjects {
			if v == nullEntry {
				indexOfNull = i + 1
				break
			}
		}

		tagsKeyMapping := tagObjects[indexOfNull].(map[string]interface{})

		indexTagUUID := uint64(tagsKeyMapping[SnippetTagsTagUUID].(plist.UID))
		indexTagTitle := uint64(tagsKeyMapping[SnippetTagsTagTitle].(plist.UID))

		uuid := tagObjects[indexTagUUID].(string)
		title := tagObjects[indexTagTitle].(string)

		result[uuid] = title
	}

	return result, nil
}

func parseSnippets(library snippetsLabLibrary) ([]model.Snippet, error) {
	filePath, err := library.snippetsFilePath()
	if err != nil {
		return []model.Snippet{}, err
	}

	// Open the directory.
	dirRead, _ := os.Open(filepath.Clean(filePath))

	// Call Readdir to get all files.
	dirFiles, _ := dirRead.Readdir(0)

	var snippets []model.Snippet
	for i := range dirFiles {
		file := dirFiles[i]

		if snippet, err := parseSnippet(fmt.Sprintf("%s/%s", filePath, file.Name())); err != nil {
			return snippets, err
		} else {
			snippets = append(snippets, snippet)
		}
	}

	return snippets, nil
}

//nolint:forcetypeassert,funlen // since we will catch any panic error and checking each statement explicitly is too much work
func parseSnippet(path string) (model.Snippet, error) {
	fileBytes, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return snippetImpl{}, err
	}

	buf := bytes.NewReader(fileBytes)
	decoder := plist.NewDecoder(buf)

	fileMap := make(map[string]interface{})
	if err := decoder.Decode(&fileMap); err != nil {
		return snippetImpl{}, err
	}

	objects := fileMap["$objects"].([]interface{})

	var indexOfNull int
	for i, v := range objects {
		if v == nullEntry {
			indexOfNull = i + 1
			break
		}
	}

	keyMapping := objects[indexOfNull].(map[string]interface{})

	uuidIndex := uint64(keyMapping[SnippetUUID].(plist.UID))
	titleIndex := uint64(keyMapping[SnippetTitle].(plist.UID))
	partsIndex := uint64(keyMapping[SnippetParts].(plist.UID))
	tagsUUIDIndex := uint64(keyMapping[SnippetTagUUIDs].(plist.UID))

	snippetUIID := objects[uuidIndex].(string)

	tagsUUIDMap := objects[tagsUUIDIndex].(map[string]interface{})
	tagsUUIDList := tagsUUIDMap["NS.objects"].([]interface{})

	var tagUUIDS []string
	for _, v := range tagsUUIDList {
		tagUUID := objects[uint64(v.(plist.UID))].(string)
		tagUUIDS = append(tagUUIDS, tagUUID)
	}

	partsMap := objects[partsIndex].(map[string]interface{})
	partsValues := partsMap["NS.objects"].([]interface{})
	partIndex0 := uint64(partsValues[0].(plist.UID))
	partMap0 := objects[partIndex0].(map[string]interface{})

	partMap0ContentIndex := uint64(partMap0[SnippetPartContent].(plist.UID))
	partMap0Content := objects[partMap0ContentIndex].(map[string]interface{})
	partMap0ContentData := partMap0Content["NS.data"].([]uint8)

	partMap0LanguageIndex := uint64(partMap0[SnippetPartLanguage].(plist.UID))
	partMap0Language := objects[partMap0LanguageIndex].(string)

	snippet := snippetImpl{
		id: snippetUIID,
		languageFunc: func() model.Language {
			return mapToLanguage(partMap0Language)
		},
		tags: tagUUIDS,
		contentFunc: func() string {
			return string(partMap0ContentData)
		},
		titleFunc: func() string {
			return objects[titleIndex].(string)
		},
	}

	return snippet, nil
}

func readPblistFile(path string) (map[string]interface{}, error) {
	fileBytes, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	buf := bytes.NewReader(fileBytes)
	decoder := plist.NewDecoder(buf)

	fileMap := make(map[string]interface{})
	if err := decoder.Decode(&fileMap); err != nil {
		return nil, err
	}

	return fileMap, nil
}

func mapToLanguage(lang string) model.Language {
	language := model.LanguageUnknown
	if lang != "" {
		if l, ok := languageMapping[lang]; ok {
			language = l
		}
	}
	return language
}
