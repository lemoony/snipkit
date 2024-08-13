package pictarinesnip

import (
	"encoding/json"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/idutil"
	"github.com/lemoony/snipkit/internal/utils/stringutil"
	"github.com/lemoony/snipkit/internal/utils/system"
	"github.com/lemoony/snipkit/internal/utils/tagutil"
)

var languageMapping = map[string]model.Language{
	"shell":    model.LanguageBash,
	"yaml":     model.LanguageYAML,
	"markdown": model.LanguageMarkdown,
}

type picatrineSnippet struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Tags    []string `json:"tags"`
	Snippet string   `json:"snippet"`
	Mode    struct {
		Name string `json:"name"`
	} `json:"mode"`
}

func parseLibrary(path string, system *system.System, tags *stringutil.StringSet) []model.Snippet {
	file, err := system.Fs.Open(path)
	if err != nil {
		panic(err)
	}

	var snippets []picatrineSnippet
	if err = json.NewDecoder(file).Decode(&snippets); err != nil {
		panic(err)
	}

	return mapToModel(snippets, tags)
}

func mapToModel(rawSnippets []picatrineSnippet, tags *stringutil.StringSet) []model.Snippet {
	var result []model.Snippet

	for i := range rawSnippets {
		raw := rawSnippets[i]

		if !tagutil.HasValidTag(*tags, raw.Tags) {
			continue
		}

		result = append(result, snippetImpl{
			id:       idutil.FormatSnippetID(raw.ID, idPrefix),
			title:    raw.Name,
			tags:     raw.Tags,
			language: mapToLanguage(raw.Mode.Name),
			content:  raw.Snippet,
		})
	}
	return result
}

// https://github.com/Pictarine/macos-snippets/blob/aeb70a4b0e04025be9b511ea5810dd41671d89e7/Snip/Model/Mode.swift
func mapToLanguage(name string) model.Language {
	if entry, ok := languageMapping[name]; ok {
		return entry
	}
	return model.LanguageUnknown
}
