package app

import (
	"encoding/json"
	"encoding/xml"

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

type ExportFormat int64

const (
	ExportFormatJSON       ExportFormat = 0
	ExportFormatPrettyJSON ExportFormat = 1
	ExportFormatXML        ExportFormat = 2
)

func (a *appImpl) ExportSnippets(fields []ExportField, format ExportFormat) string {
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

	var exportBytes []byte
	var err error

	switch format {
	case ExportFormatJSON:
		exportBytes, err = json.Marshal(export)
	case ExportFormatPrettyJSON:
		exportBytes, err = json.MarshalIndent(export, "", "    ")
	case ExportFormatXML:
		exportBytes, err = xml.MarshalIndent(export, "", "    ")
	}

	if err != nil {
		panic(err)
	}
	return string(exportBytes)
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
	ID         string          `json:"id,omitempty" xml:"id,omitempty"`
	Title      string          `json:"title,omitempty" xml:"title,omitempty"`
	Content    string          `json:"content,omitempty" xml:"content,omitempty"`
	Parameters []parameterJSON `json:"parameters,omitempty" xml:"parameters,omitempty"`
}

type parameterJSON struct {
	Key          string            `json:"key" xml:"key"`
	Name         string            `json:"name" xml:"name"`
	Type         parameterTypeJSON `json:"type" xml:"type"`
	Description  string            `json:"description,omitempty" xml:"description,omitempty"`
	DefaultValue string            `json:"defaultValue,omitempty" xml:"defaultValue,omitempty"`
	Values       []string          `json:"values,omitempty" xml:"values,omitempty"`
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
