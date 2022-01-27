package pictarinesnip

import (
	"fmt"
	"strings"

	"github.com/phuslu/log"

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

	if !manager.config.Enabled {
		log.Debug().Msg("No pictarinesnip manager because it is disabled")
		return nil, nil
	}

	return manager, nil
}

func (p Manager) Info() []model.InfoLine {
	var lines []model.InfoLine

	lines = append(lines, model.InfoLine{
		IsError: false,
		Key:     "Pictarine Snip library path",
		Value:   stringutil.StringOrDefault(p.config.LibraryPath, "not set"),
	})

	lines = append(lines, model.InfoLine{
		IsError: false,
		Key:     "Pictarine Snip tags",
		Value:   stringutil.StringOrDefault(strings.Join(p.config.IncludeTags, ","), "None"),
	})

	lines = append(lines, model.InfoLine{
		IsError: true, Key: "Pictarine Snip total number of snippets", Value: fmt.Sprintf("%d", len(p.GetSnippets())),
	})

	return lines
}

func (p *Manager) GetSnippets() []model.Snippet {
	tags := p.getValidTagUUIDs()
	return parseLibrary(p.config.LibraryPath, p.system, &tags)
}

func (p *Manager) getValidTagUUIDs() stringutil.StringSet {
	result := stringutil.StringSet{}
	for _, validTag := range p.config.IncludeTags {
		result.Add(validTag)
	}
	return result
}
