package config

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"

	"emperror.dev/errors"
	"github.com/phuslu/log"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/lemoony/snippet-kit/internal/providers/filesystem"
	"github.com/lemoony/snippet-kit/internal/providers/snippetslab"
	"github.com/lemoony/snippet-kit/internal/ui"
	"github.com/lemoony/snippet-kit/internal/ui/uimsg"
	"github.com/lemoony/snippet-kit/internal/utils"
	"github.com/lemoony/snippet-kit/internal/utils/pathutil"
)

type yamlCommentKind int

type yamlComment struct {
	value string
	kind  yamlCommentKind
}

const (
	fileModeConfig  = os.FileMode(0o600)
	yamlCommentLine = yamlCommentKind(1)
	yamlCommentHead = yamlCommentKind(2)

	yamlDefaultIndent = 2
)

func createConfigFile(system *utils.System, viper *viper.Viper, term ui.Terminal) {
	config := VersionWrapper{
		Version: "1.0.0",
		Config:  Config{},
	}

	config.Config.Style = ui.DefaultConfig()
	config.Config.Providers.SnippetsLab = snippetslab.AutoDiscoveryConfig(system)
	config.Config.Providers.FileSystem = filesystem.AutoDiscoveryConfig(system)

	data, err := serializeToYamlWithComment(config)
	if err != nil {
		panic(errors.Wrap(err, "failed to serialize config to yaml"))
	}

	configPath := viper.ConfigFileUsed()

	log.Debug().Msgf("Going to use config path %s", configPath)
	err = pathutil.CreatePath(system.Fs, configPath)
	if err != nil {
		panic(errors.Wrap(err, "failed to create path to config file"))
	}

	err = afero.WriteFile(system.Fs, configPath, data, fileModeConfig)
	if err != nil {
		panic(errors.Wrap(err, "failed to write config file"))
	}

	term.PrintMessage(uimsg.ConfigFileCreate(configPath))
}

func serializeToYamlWithComment(value interface{}) ([]byte, error) {
	// get all tag comments
	commentMap := map[string][]yamlComment{}
	traverseYamlTagComments(reflect.TypeOf(value), []string{}, &commentMap)

	// parse raw yaml string into yaml.Node
	var tree yaml.Node
	if initialBytes, err := marshalToYAML(value); err != nil {
		return nil, err
	} else if err := yaml.Unmarshal(initialBytes, &tree); err != nil {
		return nil, err
	}

	// traverse yaml tree to get a map of all node paths
	treeMap := map[string]*yaml.Node{}
	traverseYamlTree(&tree, []string{}, &treeMap)

	// set the comments
	for key, comments := range commentMap {
		for _, comment := range comments {
			switch comment.kind {
			case yamlCommentLine:
				treeMap[key].LineComment = comment.value
			case yamlCommentHead:
				treeMap[key].HeadComment = comment.value
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
		if field.Type.Kind() == reflect.Struct {
			traverseYamlTagComments(field.Type, append(path, yamlName), commentsMap)
		}
	}
}

func marshalToYAML(value interface{}) ([]byte, error) {
	buf := bytes.NewBufferString("")
	encoder := yaml.NewEncoder(buf)
	encoder.SetIndent(yamlDefaultIndent)
	if err := encoder.Encode(value); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func traverseYamlTree(node *yaml.Node, path []string, treeMap *map[string]*yaml.Node) {
	if node.Kind == yaml.DocumentNode {
		for i := range node.Content {
			traverseYamlTree(node.Content[i], path, treeMap)
		}
	} else if node.Kind == yaml.MappingNode {
		for i := range node.Content {
			if node.Content[i].Kind == yaml.MappingNode {
				if i%2 == 1 {
					traverseYamlTree(node.Content[i], append(path, node.Content[i-1].Value), treeMap)
				}
			} else if node.Content[i].Kind == yaml.ScalarNode || node.Content[i].Kind == yaml.SequenceNode {
				if i%2 == 1 {
					pathPrefix := strings.Join(path, ".")
					if len(path) > 0 {
						pathPrefix = fmt.Sprintf("%s.", pathPrefix)
					}

					pathToNode := strings.TrimSpace(fmt.Sprintf("%s%s\n", pathPrefix, node.Content[i-1].Value))
					(*treeMap)[pathToNode] = node.Content[i-1]
				}
			}
		}
	}
}
