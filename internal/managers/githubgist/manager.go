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

// TODO: Document this manager in /docs

const (
	SecretKeyPAT        = cache.SecretKey("GitHub Personal Access Token")
	SecretKeyOauthToken = cache.SecretKey("GitHub OAuth Access Token") // TODO
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

func (m *Manager) Sync(events model.SyncEventChannel) {
	log.Trace().Msg("github gist sync started")

	var lines []model.SyncLine
	events <- model.SyncEvent{Status: model.SyncStatusStarted, Lines: lines}

	currentStore := m.getStoreFromCache()
	updatedStore := &store{Version: storeVersion}
	for _, gistConfig := range m.config.Gists {
		lines = append(lines, model.SyncLine{Type: model.SyncLineTypeInfo, Value: fmt.Sprintf("Checking %s", gistConfig.URL)})

		token, err := m.authToken(gistConfig, lines, events)
		if err != nil {
			panic(err)
		}

		var currentGistStore *gistStore
		if currentStore != nil {
			currentGistStore = currentStore.getGists(gistConfig)
		}

		if s := m.getSnippetsFromAPI(gistConfig, token, currentGistStore); s != nil {
			updatedStore.Gists = append(updatedStore.Gists, *s)
		}
	}

	events <- model.SyncEvent{Status: model.SyncStatusFinished, Lines: lines}

	m.storeInCache(updatedStore)

	log.Trace().Msg("github gist sync finished")
	close(events)
}

func (m *Manager) GetSnippets() []model.Snippet {
	var result []model.Snippet

	if cacheStore := m.getStoreFromCache(); cacheStore != nil {
		for _, gstore := range cacheStore.Gists {
			gistConfig := m.config.getGistConfig(gstore.URL)
			for _, raw := range gstore.RawSnippets {
				result = append(result, parseSnippet(raw, *gistConfig))
			}
		}
	}

	return result
}

func (m *Manager) authToken(cfg GistConfig, lines []model.SyncLine, events model.SyncEventChannel) (string, error) {
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

	if token, tokenFound := m.cache.GetSecret(SecretKeyPAT, cfg.URL); tokenFound {
		if tokenOK := m.checkToken(cfg, token); tokenOK {
			return token, true
		} else {
			log.Info().Msgf("Stored token for %s is invalid. Delete it.", cfg.URL)
			m.cache.DeleteSecret(SecretKeyPAT, cfg.URL)
			lines = append(lines, model.SyncLine{Type: model.SyncLineTypeError, Value: "The current token is invalid"})
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
		if ok := m.checkToken(cfg, token); !ok {
			panic("invalid token")
		}

		m.cache.PutSecret(SecretKeyPAT, cfg.URL, token)

		return token, true
	}

	if value.Abort {
		events <- model.SyncEvent{
			Status: model.SyncStatusAborted,
			Lines:  append(lines, model.SyncLine{Type: model.SyncLineTypeInfo, Value: "Aborted"}),
		}
	}

	return "", false
}

func (m *Manager) getSnippetsFromAPI(cfg GistConfig, token string, cache *gistStore) *gistStore {
	etag := ""
	if cache != nil {
		log.Debug().Msg("cached previous store available")
		etag = cache.ETag
	}

	var snippets []rawSnippet
	resp := m.getGists(cfg, etag, token)

	if !resp.hasUpdates {
		return cache
	}

	for _, gist := range *resp.gistsResponse {
		for _, file := range gist.Files {
			id := fmt.Sprintf("%s-%s", gist.ID, file.Filename)
			fileETag := ""
			var prevRawSnippet *rawSnippet
			if cache != nil {
				for i := range cache.RawSnippets {
					if cache.RawSnippets[i].ID == id {
						prevRawSnippet = &cache.RawSnippets[i]
						fileETag = cache.RawSnippets[i].ETag
						log.Trace().Msgf("Previous etag for %s: %s", id, fileETag)
						break
					}
				}
			}

			singleRawGistResp := m.getRawGist(file.RawURL, fileETag, token)

			if !singleRawGistResp.hasUpdates {
				snippets = append(snippets, *prevRawSnippet)
			} else {
				snippets = append(snippets, rawSnippet{
					ID:          id,
					Filename:    file.Filename,
					Content:     *singleRawGistResp.rawContent,
					Pubic:       gist.Public,
					Description: gist.Description,
					Language:    file.Language,
					ETag:        singleRawGistResp.etag,
					FilesInGist: len(gist.Files),
				})
			}
		}
	}

	return &gistStore{URL: cfg.URL, ETag: resp.etag, RawSnippets: snippets}
}
