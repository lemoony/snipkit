package pet

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/afero"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/system"
)

var placeholderRegex = regexp.MustCompile(`<(.*?)>`)

type tomlSnippetsFile struct {
	Snippets []tomlSnippet
}

type tomlSnippet struct {
	Description string   `toml:"description"`
	Command     string   `toml:"command"`
	Tags        []string `toml:"tag"`
}

func parseSnippetFilePaths(s *system.System) ([]string, error) {
	configPath := filepath.Join(s.UserHome(), defaultConfigPath)
	if exists, err := afero.Exists(s.Fs, configPath); err != nil {
		return nil, err
	} else if !exists {
		return []string{}, nil
	}

	configContents := string(s.ReadFile(configPath))
	snippetFilePaths, err := parseConfigForSnippetFilePaths(configContents)
	if err != nil {
		return nil, err
	}

	return snippetFilePaths, nil
}

func parseConfigForSnippetFilePaths(configContents string) ([]string, error) {
	data := map[string]map[string]interface{}{}
	_, err := toml.Decode(configContents, &data)
	if err != nil {
		return nil, err
	}

	var paths []string
	for k := range data {
		entries := data[k]
		if snippetFile, ok := entries["snippetfile"]; ok {
			if snippetFileStr, isString := snippetFile.(string); isString {
				paths = append(paths, snippetFileStr)
			}
		}
	}

	return paths, nil
}

func parseSnippetsFromTOML(contents string) []model.Snippet {
	var snippetsFile tomlSnippetsFile
	_, err := toml.Decode(contents, &snippetsFile)
	if err != nil {
		panic(err)
	}

	result := make([]model.Snippet, len(snippetsFile.Snippets))
	for i := range snippetsFile.Snippets {
		result[i] = mapToSnippet(snippetsFile.Snippets[i])
	}
	return result
}

func mapToSnippet(raw tomlSnippet) model.Snippet {
	return &snippetImpl{
		id:       "not_used",
		title:    raw.Description,
		content:  raw.Command,
		tags:     raw.Tags,
		language: model.LanguageBash,
	}
}

func parseParameters(command string) []model.Parameter {
	const expectedMatches = 2

	testRegex := regexp.MustCompile(`\|\_[^\|\|]+\_\|`)

	var result []model.Parameter
	matches := placeholderRegex.FindAllStringSubmatch(command, -1)
	for i := range matches {
		if len(matches[i]) >= expectedMatches {
			split := strings.SplitN(matches[i][1], "=", 2)
			key := strings.TrimSpace(split[0])
			defaultValue := ""
			var values []string
			if len(split) > 1 {
				if multipleDefValues := testRegex.FindAllStringSubmatch(split[1], -1); multipleDefValues != nil {
					println(multipleDefValues)
					for _, val := range multipleDefValues {
						values = append(values, strings.Trim(val[0], "|_"))
					}
				} else {
					defaultValue = strings.TrimSpace(split[1])
				}
			}
			result = append(result, model.Parameter{Key: key, DefaultValue: defaultValue, Values: values})
		}
	}
	return result
}

func formatContent(command string, values []string) string {
	if len(values) == 0 {
		return command
	}

	args := make([]interface{}, len(values))
	for i := range values {
		args[i] = values[i]
	}
	return fmt.Sprintf(placeholderRegex.ReplaceAllString(command, "%s"), args...)
}
