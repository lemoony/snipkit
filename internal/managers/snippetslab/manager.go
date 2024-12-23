package snippetslab

import (
	"fmt"
	"strings"

	"emperror.dev/errors"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/stringutil"
	"github.com/lemoony/snipkit/internal/utils/system"
	"github.com/lemoony/snipkit/internal/utils/tagutil"
)

type Manager struct {
	system *system.System
	config Config
}

// Option configures a Manager.
type Option interface {
	apply(p *Manager)
}

// optionFunc wraps a func so that it satisfies the Option interface.
type optionFunc func(provider *Manager)

func (f optionFunc) apply(provider *Manager) {
	f(provider)
}

// WithSystem sets the utils.System instance to be used by Manager.
func WithSystem(system *system.System) Option {
	return optionFunc(func(p *Manager) {
		p.system = system
	})
}

func WithConfig(config Config) Option {
	return optionFunc(func(p *Manager) {
		p.config = config
	})
}

func NewManager(options ...Option) (*Manager, error) {
	manager := &Manager{}
	for _, o := range options {
		o.apply(manager)
	}
	return manager, nil
}

func (m Manager) libraryPath() snippetsLabLibrary {
	return snippetsLabLibrary(m.config.LibraryPath)
}

func (m Manager) Key() model.ManagerKey {
	return Key
}

func (m Manager) Info() []model.InfoLine {
	var lines []model.InfoLine

	if path, err := findPreferencesPath(m.system); err != nil {
		lines = append(lines, model.InfoLine{IsError: true, Key: "SnippetsLab preferences path", Value: err.Error()})
	} else {
		lines = append(lines, model.InfoLine{
			IsError: true, Key: "SnippetsLab preferences path", Value: path,
		})
	}

	lines = append(lines, model.InfoLine{
		IsError: false, Key: "SnippetsLab library path", Value: m.config.LibraryPath,
	})

	lines = append(lines, model.InfoLine{
		IsError: false,
		Key:     "SnippetsLab tags",
		Value:   stringutil.StringOrDefault(strings.Join(m.config.IncludeTags, ","), "None"),
	})

	lines = append(lines, model.InfoLine{
		IsError: true, Key: "SnippetsLab total number of snippets", Value: fmt.Sprintf("%d", len(m.GetSnippets())),
	})

	return lines
}

func (m *Manager) Sync(model.SyncEventChannel) {
	// do nothing
}

func (m Manager) SaveAssistantSnippet(snippetTitle string, filename string, contents []byte) {
	panic(errors.New("Not implemented"))
}

func (m *Manager) GetSnippets() []model.Snippet {
	validTagUUIDs, hasTags := m.getValidTagUUIDs()

	snippets, err := parseSnippets(m.libraryPath())
	if err != nil {
		panic(err)
	}

	if hasTags {
		var result []model.Snippet
		if len(validTagUUIDs) > 0 {
			for _, snippet := range snippets {
				if tagutil.HasValidTag(validTagUUIDs, snippet.GetTags()) {
					result = append(result, snippet)
				}
			}
		}
		return result
	}

	return snippets
}

func (m *Manager) getValidTagUUIDs() (stringutil.StringSet, bool) {
	tags, err := parseTags(m.libraryPath())
	if err != nil {
		panic(err)
	}

	result := stringutil.StringSet{}
	for _, validTag := range m.config.IncludeTags {
		for tagKey, tagValue := range tags {
			if strings.Compare(tagValue, validTag) == 0 {
				result.Add(tagKey)
			}
		}
	}

	return result, len(m.config.IncludeTags) > 0
}
