package githubgist

import (
	"fmt"
	"regexp"
	"time"

	"github.com/phuslu/log"

	"github.com/lemoony/snipkit/internal/cache"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/system"
)

type Manager struct {
	system      *system.System
	config      Config
	suffixRegex []*regexp.Regexp //nolint:structcheck,unused // ignore for now since not used yet
	cache       cache.Cache
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

func WithCache(cache cache.Cache) Option {
	return optionFunc(func(p *Manager) {
		p.cache = cache
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

func (m Manager) Key() model.ManagerKey {
	return Key
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

func (m *Manager) Sync(events model.SyncEventChannel) bool {
	log.Trace().Msg("github gist sync started")

	var lines []model.SyncLine
	events <- model.SyncEvent{Status: model.SyncStatusStarted, Lines: lines}

	time.Sleep(time.Second * 1)

	for _, g := range m.config.Gists {
		lines = append(lines, model.SyncLine{Type: model.SyncLineTypeInfo, Value: fmt.Sprintf("Checking %s", g.URL)})
		if g.AuthenticationMethod == AuthMethodToken {
			if _, ok := m.requestAuthToken(lines, events); !ok {
				return false
			}
		}
	}

	time.Sleep(time.Second * 2) //nolint:gomnd //ignore for now

	lines = append(lines, model.SyncLine{
		Type:  model.SyncLineTypeSuccess,
		Value: "Input successful. Stored token in keychain.",
	})

	events <- model.SyncEvent{
		Status: model.SyncStatusStarted,
		Lines:  lines,
		Login:  nil,
	}

	time.Sleep(time.Second * 2) //nolint:gomnd //ignore for now
	snippets := m.getSnippetsFromAPI()
	events <- model.SyncEvent{Status: model.SyncStatusFinished, Lines: lines}
	close(events)
	log.Trace().Msg("github gist sync finished")
	fmt.Println(snippets)

	return true
}

func (m *Manager) GetSnippets() []model.Snippet {
	var result []model.Snippet
	return result
}

func (m *Manager) getSnippetsFromAPI() []model.Snippet {
	var snippets []model.Snippet
	for _, cfg := range m.config.Gists {
		resp, err := m.getGists(cfg)
		if err != nil {
			fmt.Println(err)
		}

		for _, gist := range resp {
			for _, file := range gist.Files {
				rawResp, err := m.getRawGist(file.RawURL)
				if err != nil {
					fmt.Println(err)
				}

				snippets = append(snippets, model.Snippet{
					UUID: fmt.Sprintf("%s-%s", gist.ID, file.Filename),
					TitleFunc: func() string {
						return file.Filename
					},
					ContentFunc: func() string {
						return string(rawResp)
					},
				})
			}
		}
	}
	return snippets
}

func (m *Manager) requestAuthToken(lines []model.SyncLine, events model.SyncEventChannel) (string, bool) {
	contChannel := make(chan model.SyncInputResult)

	events <- model.SyncEvent{
		Status: model.SyncStatusStarted,
		Lines:  lines,
		Login: &model.SyncInput{
			Content:     "You need to login into GitHub.\nYou have not yet provided an Access Token..",
			Placeholder: "Access token",
			Type:        model.SyncLoginTypeText,
			Input:       contChannel,
		},
	}

	//nolint:ifshort //TODO: refactor at a later point
	value := <-contChannel

	events <- model.SyncEvent{Status: model.SyncStatusStarted, Lines: lines}

	if value.Text != "" {
		return value.Text, true
	}

	if value.Abort {
		events <- model.SyncEvent{
			Status: model.SyncStatusStarted,
			Lines:  append(lines, model.SyncLine{Type: model.SyncLineTypeInfo, Value: "Aborted"}),
		}
	}

	return "", false
}
