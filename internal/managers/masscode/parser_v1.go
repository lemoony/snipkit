package masscode

import (
	"encoding/json"
	"path/filepath"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/system"
)

func parseDBFileV1(sys *system.System, massCodePath string) []model.Snippet {
	var result []model.Snippet

	tagMap := parseRawTagMapV1(sys, filepath.Join(massCodePath, v1TagsFile))
	snippetsMap := parseRawSnippetsV1(sys, filepath.Join(massCodePath, v1SnippetsFile))

	for _, raw := range snippetsMap {
		result = append(result, &snippetImpl{
			id:       raw.ID,
			title:    raw.Name,
			tags:     toTagNames(raw.Tags, tagMap),
			content:  raw.Content[0].Value,
			language: mapLanguage(raw.Content[0].Language),
		})
	}
	return result
}

func parseRawSnippetsV1(sys *system.System, path string) map[string]rawSnippet {
	snippets := map[string]rawSnippet{}

	file, err := sys.Fs.Open(path)
	if err != nil {
		panic(err)
	}

	dc := json.NewDecoder(file)

	var snippet rawSnippet
	for err = dc.Decode(&snippet); err == nil; err = dc.Decode(&snippet) {
		if snippet.Deleted || snippet.IsInTrash {
			delete(snippets, snippet.ID)
		} else {
			snippets[snippet.ID] = snippet
		}
		snippet = rawSnippet{}
	}

	return snippets
}

func parseRawTagMapV1(sys *system.System, path string) map[string]string {
	tags := map[string]string{}

	file, err := sys.Fs.Open(path)
	if err != nil {
		panic(err)
	}

	dc := json.NewDecoder(file)

	var tag rawTag
	for err = dc.Decode(&tag); err == nil; err = dc.Decode(&tag) {
		if tag.Deleted {
			delete(tags, tag.ID)
		} else {
			tags[tag.ID] = tag.Name
		}
	}

	return tags
}
