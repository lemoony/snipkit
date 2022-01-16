package config

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"emperror.dev/errors"
	"github.com/phuslu/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/lemoony/snipkit/internal/managers/fslibrary"
	"github.com/lemoony/snipkit/internal/managers/pictarinesnip"
	"github.com/lemoony/snipkit/internal/managers/snippetslab"
	"github.com/lemoony/snipkit/internal/ui"
	"github.com/lemoony/snipkit/internal/ui/uimsg"
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
)

func createConfigFile(system *system.System, viper *viper.Viper, term ui.Terminal) {
	config := VersionWrapper{
		Version: "1.0.0",
		Config:  Config{},
	}

	config.Config.Style = ui.DefaultConfig()
	config.Config.Manager.SnippetsLab = snippetslab.AutoDiscoveryConfig(system)
	config.Config.Manager.PictarineSnip = pictarinesnip.AutoDiscoveryConfig(system)
	config.Config.Manager.FsLibrary = fslibrary.AutoDiscoveryConfig(system)

	data := serializeToYamlWithComment(config)

	configPath := viper.ConfigFileUsed()

	log.Debug().Msgf("Going to use config path %s", configPath)
	system.CreatePath(configPath)
	system.WriteFile(configPath, data)

	term.PrintMessage(uimsg.ConfigFileCreated(configPath))
}

func serializeToYamlWithComment(value interface{}) []byte {
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
