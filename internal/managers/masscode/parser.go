package masscode

import (
	"encoding/json"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/system"
)

var languageMapping = map[string]model.Language{
	"shell":    model.LanguageBash,
	"yaml":     model.LanguageYAML,
	"markdown": model.LanguageMarkdown,
	"toml":     model.LanguageTOML,
}

type v2DbFile struct {
	Snippets []rawSnippet `json:"Snippets"`
	Tags     []rawTag     `json:"Tags"`
}

type rawTag struct {
	Name string `json:"name"`
	ID   string `json:"_id"`
}

type rawSnippet struct {
	Name    string   `json:"name"`
	ID      string   `json:"_id"`
	TagIDs  []string `json:"tagIds"`
	Content []struct {
		Language string `json:"language"`
		Value    string `json:"value"`
	}
}

func parseDBFileV2(sys *system.System, path string) []model.Snippet {
	var result []model.Snippet

	var dbFile v2DbFile
	contents := sys.ReadFile(path)
	if err := json.Unmarshal(contents, &dbFile); err != nil {
		panic(err)
	}

	tagMap := toTagMap(dbFile.Tags)

	for _, raw := range dbFile.Snippets {
		var tags []string
		tagIds := raw.TagIDs
		for t := range raw.TagIDs {
			if tag, ok := tagMap[tagIds[t]]; ok {
				tags = append(tags, tag)
			}
		}

		result = append(result, &snippetImpl{
			id:       raw.ID,
			title:    raw.Name,
			tags:     tags,
			content:  raw.Content[0].Value,
			language: mapLanguage(raw.Content[0].Language),
		})
	}

	return result
}

func toTagMap(tags []rawTag) map[string]string {
	result := map[string]string{}
	for i := range tags {
		result[tags[i].ID] = tags[i].Name
	}
	return result
}

func mapLanguage(value string) model.Language {
	if l, ok := languageMapping[value]; ok {
		return l
	}
	return model.LanguageText
}
