package masscode

import (
	"fmt"
	"path/filepath"

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

	lines = append(lines, model.InfoLine{Key: "MassCode enabled", Value: fmt.Sprintf("%v", m.config.Enabled)})
	lines = append(lines, model.InfoLine{Key: "MassCode path", Value: m.config.MassCodeHome})

	lines = append(lines, model.InfoLine{
		Key:   "MassCode total number of Snippets",
		Value: fmt.Sprintf("%d", len(m.GetSnippets())),
	})

	return lines
}

func (m *Manager) GetSnippets() []model.Snippet {
	var result []model.Snippet
	validTags := stringutil.NewStringSet(m.config.IncludeTags)

	var snippets []model.Snippet
	if m.config.Version == version1 {
		// TODO
	} else if m.config.Version == version2 {
		snippets = parseDBFileV2(m.system, filepath.Join(m.system.UserHome(), defaultMassCodeHomePath, v2DatabaseFile))
	}
	for _, snippet := range snippets {
		if tagutil.HasValidTag(validTags, snippet.GetTags()) {
			result = append(result, snippet)
		}
	}
	return result
}
