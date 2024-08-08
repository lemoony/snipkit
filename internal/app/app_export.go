package app

import (
	"encoding/json"

	"golang.org/x/exp/slices"

	"github.com/lemoony/snipkit/internal/model"
)

type ExportField int64

const (
	ExportFieldID         ExportField = 0
	ExportFieldTitle      ExportField = 1
	ExportFieldContent    ExportField = 2
	ExportFieldParameters ExportField = 3
)

func (a *appImpl) ExportSnippets(fields []ExportField) string {
	snippets := a.getAllSnippets()
	if len(snippets) == 0 {
		panic(ErrNoSnippetsAvailable)
	}

	var snippetJSONList []snippetJSON
	for _, snippet := range snippets {
		snippetJSONList = append(snippetJSONList, convertSnippetToJSON(snippet, fields))
	}

	export := exportJSON{
		Snippets: snippetJSONList,
	}

	jsonData, err := json.Marshal(export)
	if err != nil {
		panic(err)
	}

	return string(jsonData)
}

func convertSnippetToJSON(snippet model.Snippet, fields []ExportField) snippetJSON {
	return snippetJSON{
		ID: func() string {
			if slices.Contains(fields, ExportFieldID) {
				return snippet.GetID()
			}
			return ""
		}(),
		Title: func() string {
			if slices.Contains(fields, ExportFieldTitle) {
				return snippet.GetTitle()
			}
			return ""
		}(),
		Content: func() string {
			if slices.Contains(fields, ExportFieldContent) {
				return snippet.GetContent()
			}
			return ""
		}(),
		Parameters: func() []parameterJSON {
			if slices.Contains(fields, ExportFieldParameters) {
				return convertParametersToJSON(snippet.GetParameters())
			}
			return []parameterJSON{}
		}(),
	}
}

func convertParametersToJSON(parameters []model.Parameter) []parameterJSON {
	result := make([]parameterJSON, len(parameters))
	for i, v := range parameters {
		result[i] = parameterJSON{
			Key:          v.Key,
			Name:         v.Name,
			Description:  v.Description,
			DefaultValue: v.DefaultValue,
			Type:         parameterTypeMap[v.Type],
			Values:       v.Values,
		}
	}
	return result
}

type exportJSON struct {
	Snippets []snippetJSON `json:"snippets"`
}

type snippetJSON struct {
	ID         string          `json:"id,omitempty"`
	Title      string          `json:"title,omitempty"`
	Content    string          `json:"content,omitempty"`
	Parameters []parameterJSON `json:"parameters,omitempty"`
}

type parameterJSON struct {
	Key          string
	Name         string
	Type         parameterTypeJSON
	Description  string
	DefaultValue string
	Values       []string
}

type parameterTypeJSON string

const (
	parameterTypeValue    parameterTypeJSON = "VALUE"
	parameterTypePath     parameterTypeJSON = "PATH"
	parameterTypePassword parameterTypeJSON = "PASSWORD"
)

var parameterTypeMap = map[model.ParameterType]parameterTypeJSON{
	model.ParameterTypeValue:    parameterTypeValue,
	model.ParameterTypePath:     parameterTypePath,
	model.ParameterTypePassword: parameterTypePassword,
}
