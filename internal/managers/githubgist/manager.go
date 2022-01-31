package githubgist

import (
	"fmt"
	"regexp"
	"time"

	"github.com/phuslu/log"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/system"
)

type Manager struct {
	system      *system.System
	config      Config
	suffixRegex []*regexp.Regexp //nolint:structcheck,unused // ignore for now since not used yet
}

// Option configures a Manager.
type Option interface {
	apply(p *Manager)
}

// optionFunc wraps a func so that it satisfies the Option interface.
type optionFunc func(manager *Manager)

func (f optionFunc) apply(manager *Manager) {
	f(manager)
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
		log.Debug().Msg("No github gist manager because it is disabled")
		return nil, nil
	}

	return manager, nil
}

func (m Manager) Info() []model.InfoLine {
	var lines []model.InfoLine

	lines = append(lines, model.InfoLine{
		IsError: false,
		Key:     "GitHub Gist enabled",
		Value:   fmt.Sprintf("%v", m.config.Enabled),
	})

	lines = append(lines, model.InfoLine{
		IsError: false,
		Key:     "GitHub Gist number of URLs",
		Value:   fmt.Sprintf("%d", len(m.config.Gists)),
	})

	lines = append(lines, model.InfoLine{
		IsError: false, Key: "GitHub Gist total number of snippets", Value: fmt.Sprintf("%d", len(m.GetSnippets())),
	})

	return lines
}

func (m *Manager) Sync() model.SyncResult {
	time.Sleep(time.Second * 3) //nolint:gomnd //ignore for now
	return model.SyncResult{
		Added:   0,
		Updated: 0,
		Deleted: 0,
	}
}

func (m *Manager) GetSnippets() []model.Snippet {
	var result []model.Snippet
	return result
}
