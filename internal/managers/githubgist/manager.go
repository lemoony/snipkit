package githubgist

import (
	"fmt"
	"regexp"

	"emperror.dev/errors"
	"github.com/phuslu/log"

	"github.com/lemoony/snipkit/internal/cache"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/system"
)

const (
	SecretKeyPAT        = cache.SecretKey("GitHub Personal Access Token")
	SecretKeyOauthToken = cache.SecretKey("GitHub OAuth Access Token")
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

	var snippets []model.Snippet
	for _, gistConfig := range m.config.Gists {
		lines = append(lines, model.SyncLine{Type: model.SyncLineTypeInfo, Value: fmt.Sprintf("Checking %s", gistConfig.gistURL())})

		token, err := m.authToken(gistConfig, lines, events)
		if err != nil {
			panic(err)
		}

		if s, err2 := m.getSnippetsFromAPI(gistConfig, token); err2 != nil {
			panic(err2)
		} else {
			snippets = append(snippets, s...)
		}
	}

	events <- model.SyncEvent{Status: model.SyncStatusFinished, Lines: lines}
	close(events)
	log.Trace().Msg("github gist sync finished")
	log.Trace().Msgf("number of retrieved snippets: %d", len(snippets))
	return true
}

func (m *Manager) GetSnippets() []model.Snippet {
	var result []model.Snippet
	return result
}

func (m *Manager) authToken(cfg GistConfig, lines []model.SyncLine, events model.SyncEventChannel) (string, error) {
	if cfg.AuthenticationMethod == AuthMethodNone {
		return "", nil
	}

	switch cfg.AuthenticationMethod {
	case AuthMethodNone:
		return "", nil
	case AuthMethodToken:
		if token, ok := m.requestAuthToken(cfg, lines, events); !ok {
			return "", errors.New("no auth token")
		} else {
			return token, nil
		}
	}

	panic("TODO: oauth not supported yet")
}

func (m *Manager) requestAuthToken(cfg GistConfig, lines []model.SyncLine, events model.SyncEventChannel) (string, bool) {
	contChannel := make(chan model.SyncInputResult)

	account := fmt.Sprintf("%s/%s", cfg.Host, cfg.Username)
	if token, tokenFound := m.cache.GetSecret(SecretKeyPAT, account); tokenFound {
		tokenOK, tokenErr := m.checkToken(cfg, token)
		switch {
		case tokenErr != nil:
			panic(tokenErr)
		case !tokenOK:
			log.Info().Msgf("Stored token for %s is invalid. Delete it.", account)
			m.cache.DeleteSecret(SecretKeyPAT, account)
			lines = append(lines, model.SyncLine{Type: model.SyncLineTypeError, Value: "The current token is invalid"})
		case tokenOK:
			return token, true
		}
	}

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

	value := <-contChannel

	events <- model.SyncEvent{Status: model.SyncStatusStarted, Lines: lines}

	if token := value.Text; token != "" {
		if ok, err := m.checkToken(cfg, token); err != nil {
			panic(err)
		} else if !ok {
			panic("invalid token")
		}

		m.cache.PutSecret(SecretKeyPAT, account, token)

		return token, true
	}

	if value.Abort {
		events <- model.SyncEvent{
			Status: model.SyncStatusStarted,
			Lines:  append(lines, model.SyncLine{Type: model.SyncLineTypeInfo, Value: "Aborted"}),
		}
	}

	return "", false
}

func (m *Manager) getSnippetsFromAPI(cfg GistConfig, token string) ([]model.Snippet, error) {
	var snippets []model.Snippet
	resp, err := m.getGists(cfg, token)
	if err != nil && errors.Is(err, errAuth) {
		return nil, err
	} else if err != nil {
		panic(err)
	}

	for _, gist := range resp {
		for _, file := range gist.Files {
			rawResp, err := m.getRawGist(file.RawURL)
			if err != nil {
				panic(err)
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
	return snippets, nil
}
