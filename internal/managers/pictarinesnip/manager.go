package pictarinesnip

import (
	"fmt"
	"strings"

	"emperror.dev/errors"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/stringutil"
	"github.com/lemoony/snipkit/internal/utils/system"
)

type Manager struct {
	system *system.System
	config Config
}

// Option configures a Manager.
type Option interface {
	apply(m *Manager)
}

// optionFunc wraps a func so that it satisfies the Option interface.
type optionFunc func(m *Manager)

func (f optionFunc) apply(m *Manager) {
	f(m)
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

func (m Manager) Key() model.ManagerKey {
	return Key
}

func (m *Manager) Sync(model.SyncEventChannel) {
	// do nothing
}

func (m Manager) Info() []model.InfoLine {
	var lines []model.InfoLine

	lines = append(lines, model.InfoLine{
		IsError: false,
		Key:     "Pictarine Snip library path",
		Value:   stringutil.StringOrDefault(m.config.LibraryPath, "not set"),
	})

	lines = append(lines, model.InfoLine{
		IsError: false,
		Key:     "Pictarine Snip tags",
		Value:   stringutil.StringOrDefault(strings.Join(m.config.IncludeTags, ","), "None"),
	})

	lines = append(lines, model.InfoLine{
		IsError: true, Key: "Pictarine Snip total number of snippets", Value: fmt.Sprintf("%d", len(m.GetSnippets())),
	})

	return lines
}

func (m *Manager) GetSnippets() []model.Snippet {
	tags := m.getValidTagUUIDs()
	return parseLibrary(m.config.LibraryPath, m.system, &tags)
}

func (m Manager) SaveAssistantSnippet(snippetTitle string, filename string, contents []byte) {
	panic(errors.New("Not implemented"))
}

func (m *Manager) getValidTagUUIDs() stringutil.StringSet {
	result := stringutil.StringSet{}
	for _, validTag := range m.config.IncludeTags {
		result.Add(validTag)
	}
	return result
}
