package masscode

import (
	"encoding/json"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/idutil"
	"github.com/lemoony/snipkit/internal/utils/system"
)

var languageMapping = map[string]model.Language{
	"shell":    model.LanguageBash,
	"yaml":     model.LanguageYAML,
	"markdown": model.LanguageMarkdown,
	"toml":     model.LanguageTOML,
}

type rawTag struct {
	ID      string `json:"_id"`
	Deleted bool   `json:"$$deleted"`
	Name    string `json:"name"`
}

type rawSnippet struct {
	ID        string   `json:"_id"`
	Deleted   bool     `json:"$$deleted"`
	Name      string   `json:"name"`
	TagIDs    []string `json:"tagIds"` // used for v2.
	Tags      []string `json:"tags"`   // used for v1.
	IsInTrash bool     `json:"isDeleted"`
	Content   []struct {
		Language string `json:"language"`
		Value    string `json:"value"`
	}
}

func parseDBFileV2(sys *system.System, path string) []model.Snippet {
	var result []model.Snippet

	type v2DbFile struct {
		Snippets []rawSnippet `json:"Snippets"`
		Tags     []rawTag     `json:"Tags"`
	}

	var dbFile v2DbFile
	contents := sys.ReadFile(path)
	if err := json.Unmarshal(contents, &dbFile); err != nil {
		panic(err)
	}

	tagMap := toTagMapV2(dbFile.Tags)

	for _, raw := range dbFile.Snippets {
		result = append(result, &snippetImpl{
			id:       idutil.FormatSnippetID(raw.ID, idPrefix),
			title:    raw.Name,
			tags:     toTagNames(raw.TagIDs, tagMap),
			content:  raw.Content[0].Value,
			language: mapLanguage(raw.Content[0].Language),
		})
	}

	return result
}

func toTagMapV2(tags []rawTag) map[string]string {
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

func toTagNames(tagIDs []string, tagMap map[string]string) []string {
	var tags []string
	tagIds := tagIDs
	for t := range tagIDs {
		if tag, ok := tagMap[tagIds[t]]; ok {
			tags = append(tags, tag)
		}
	}
	return tags
}
