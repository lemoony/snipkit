package config

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"emperror.dev/errors"
	"github.com/phuslu/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/lemoony/snipkit/internal/config/migrations"
	"github.com/lemoony/snipkit/internal/ui"
	"github.com/lemoony/snipkit/internal/utils/system"
)

type yamlCommentKind int

type yamlComment struct {
	value string
	kind  yamlCommentKind
}

const (
	yamlCommentLine = yamlCommentKind(1)
	yamlCommentHead = yamlCommentKind(2)

	yamlDefaultIndent = 2

	Version = migrations.Latest
)

var sliceIndexRegex = regexp.MustCompile(`\[\d]`)

func wrap(config Config) VersionWrapper {
	return VersionWrapper{
		Version: Version,
		Config:  config,
	}
}

func createConfigFile(system *system.System, viper *viper.Viper) {
	config := wrap(defaultConfig())
	data := SerializeToYamlWithComment(config)

	configPath := viper.ConfigFileUsed()
	log.Debug().Msgf("Going to use config path %s", configPath)
	system.CreatePath(configPath)
	system.WriteFile(configPath, data)
}

func defaultConfig() Config {
	return Config{
		Style: ui.DefaultConfig(),
		Script: ScriptConfig{
			ParameterMode:  ParameterModeSet,
			RemoveComments: false,
		},
	}
}

func SerializeToYamlWithComment(value interface{}) []byte {
	// get all tag comments
	commentMap := map[string][]yamlComment{}
	traverseYamlTagComments(reflect.TypeOf(value), []string{}, &commentMap)

	// parse raw yaml string into yaml.Node
	var tree yaml.Node
	if err := yaml.Unmarshal(marshalToYAML(value), &tree); err != nil {
		panic(errors.Wrap(err, "failed to unmarshal yaml"))
	}

	// traverse yaml tree to get a map of all node paths
	treeMap := map[string]*yaml.Node{}
	traverseYamlTree(&tree, []string{}, &treeMap)

	findComment := func(key string, commentMap map[string][]yamlComment) ([]yamlComment, bool) {
		if r, ok := commentMap[sliceIndexRegex.ReplaceAllString(key, "[*]")]; ok {
			return r, ok
		}

		return nil, false
	}

	for key := range treeMap {
		if comments, ok := findComment(key, commentMap); ok {
			for _, comment := range comments {
				switch comment.kind {
				case yamlCommentLine:
					treeMap[key].LineComment = comment.value
				case yamlCommentHead:
					treeMap[key].HeadComment = comment.value
				}
			}
		}
	}

	return marshalToYAML(&tree)
}

func traverseYamlTagComments(t reflect.Type, path []string, commentsMap *map[string][]yamlComment) {
	pathPrefix := strings.Join(path, ".")
	if len(path) > 0 {
		pathPrefix = fmt.Sprintf("%s.", pathPrefix)
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		yamlName := field.Tag.Get("yaml")
		if yamlName == "" {
			continue
		} else if splits := strings.Split(yamlName, ","); len(splits) > 1 {
			yamlName = splits[0]
		}
		nodePath := strings.TrimSpace(fmt.Sprintf("%s%s\n", pathPrefix, yamlName))
		commentsList := (*commentsMap)[nodePath]
		if c := field.Tag.Get("line_comment"); c != "" {
			commentsList = append(commentsList, yamlComment{value: c, kind: yamlCommentLine})
			(*commentsMap)[nodePath] = commentsList
		}
		if c := field.Tag.Get("head_comment"); c != "" {
			(*commentsMap)[nodePath] = append(commentsList, yamlComment{value: c, kind: yamlCommentHead})
		}

		switch field.Type.Kind() {
		case reflect.Struct:
			traverseYamlTagComments(field.Type, append(path, yamlName), commentsMap)
		case reflect.Ptr:
			if field.Type.Elem().Kind() == reflect.Struct {
				traverseYamlTagComments(field.Type.Elem(), append(path, yamlName), commentsMap)
			}
		case reflect.Slice:
			if field.Type.Elem().Kind() == reflect.Struct {
				traverseYamlTagComments(field.Type.Elem(), append(path, yamlName+"[*]"), commentsMap)
			}
		}
	}
}

func marshalToYAML(value interface{}) []byte {
	buf := bytes.NewBufferString("")
	encoder := yaml.NewEncoder(buf)
	encoder.SetIndent(yamlDefaultIndent)
	if err := encoder.Encode(value); err != nil {
		panic(errors.Wrap(err, "failed to marshal to yaml"))
	}
	return buf.Bytes()
}

func traverseYamlTree(node *yaml.Node, path []string, treeMap *map[string]*yaml.Node) {
	switch node.Kind {
	case yaml.DocumentNode:
		for i := range node.Content {
			traverseYamlTree(node.Content[i], path, treeMap)
		}
	case yaml.MappingNode:
		handleMappingNode(node, path, treeMap)
	}
}

func handleMappingNode(node *yaml.Node, path []string, treeMap *map[string]*yaml.Node) {
	for i := range node.Content {
		if i%2 == 0 {
			continue // Skip keys, only process values
		}

		valueNode := node.Content[i]
		keyNode := node.Content[i-1]

		switch valueNode.Kind {
		case yaml.MappingNode:
			traverseYamlTree(valueNode, append(path, keyNode.Value), treeMap)
		case yaml.ScalarNode, yaml.SequenceNode:
			processScalarOrSequence(valueNode, keyNode, path, treeMap)
		}
	}
}

func processScalarOrSequence(valueNode, keyNode *yaml.Node, path []string, treeMap *map[string]*yaml.Node) {
	pathPrefix := strings.Join(path, ".")
	if len(path) > 0 {
		pathPrefix += "."
	}

	if valueNode.Kind == yaml.SequenceNode {
		for idx, v := range valueNode.Content {
			traverseYamlTree(v, append(path, fmt.Sprintf("%s[%d]", keyNode.Value, idx)), treeMap)
		}
	}

	(*treeMap)[strings.TrimSpace(fmt.Sprintf("%s%s\n", pathPrefix, keyNode.Value))] = keyNode
}
