package githubgist

import (
	"fmt"
	"regexp"

	"emperror.dev/errors"
	"github.com/cli/oauth"
	"github.com/phuslu/log"

	"github.com/lemoony/snipkit/internal/cache"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui/uimsg"
	"github.com/lemoony/snipkit/internal/utils/stringutil"
	"github.com/lemoony/snipkit/internal/utils/system"
	"github.com/lemoony/snipkit/internal/utils/tagutil"
)

const (
	secretKeyAccessToken = cache.SecretKey("GitHub Access Token")
	defaultOAuthClientID = "26b4225d7ae7b3961624"
)

var errAbort = errors.New("abort")

type Manager struct {
	system      *system.System
	config      Config
	suffixRegex []*regexp.Regexp //nolint:unused // ignore for now since not used yet
	cache       cache.Cache
	browseURL   func(s string) error
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

func (m *Manager) GetSnippets() []model.Snippet {
	var result []model.Snippet

	if cacheStore := m.getStoreFromCache(); cacheStore != nil {
		for _, gstore := range cacheStore.Gists {
			gistConfig := m.config.getGistConfig(gstore.URL)
			validTags := stringutil.NewStringSet(gistConfig.IncludeTags)
			for _, raw := range gstore.RawSnippets {
				snippet := parseSnippet(raw, *gistConfig)
				if tagutil.HasValidTag(validTags, snippet.GetTags()) {
					result = append(result, parseSnippet(raw, *gistConfig))
				}
			}
		}
	}

	return result
}

func (m *Manager) Sync(events model.SyncEventChannel) {
	var lines []model.SyncLine
	log.Trace().Msg("github gist sync started")

	defer func() {
		if panicValue := recover(); panicValue != nil {
			err := errors.Errorf("Sync failed: %s", panicValue)
			log.Error().Err(err).Msg("Sync failed")
			events <- model.SyncEvent{
				Status: model.SyncStatusAborted,
				Lines:  append(lines, model.SyncLine{Type: model.SyncLineTypeError, Value: err.Error()}),
			}
		}
	}()

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
}

func (m *Manager) authToken(cfg GistConfig, lines []model.SyncLine, events model.SyncEventChannel) (string, error) {
	switch cfg.AuthenticationMethod {
	case AuthMethodNone:
		return "", nil
	case AuthMethodPAT:
		if token, err := m.requestPAT(cfg, lines, events); err != nil {
			return "", err
		} else {
			return token, nil
		}
	case AuthMethodOAuthDeviceFlow:
		if token, err := m.requestOAuthToken(cfg, lines, events); err != nil {
			return "", err
		} else {
			return token, nil
		}
	}

	panic(errors.Errorf("unsupported auth method: %s", cfg.AuthenticationMethod))
}

func (m *Manager) requestPAT(cfg GistConfig, lines []model.SyncLine, events model.SyncEventChannel) (string, error) {
	contChannel := make(chan model.SyncInputResult)

	if token, tokenFound := m.cache.GetSecret(secretKeyAccessToken, cfg.URL); tokenFound {
		if tokenOK := m.checkToken(cfg, token); tokenOK {
			return token, nil
		} else {
			log.Info().Msgf("Stored token for %s is invalid. Delete it.", cfg.URL)
			m.cache.DeleteSecret(secretKeyAccessToken, cfg.URL)
			lines = append(lines, model.SyncLine{Type: model.SyncLineTypeError, Value: "The current token is invalid"})
		}
	}

	events <- model.SyncEvent{
		Status: model.SyncStatusStarted,
		Lines:  lines,
		Login: &model.SyncInput{
			Content:     fmt.Sprintf("Please provide an access token for %s with scope 'gist'...", cfg.URL),
			Placeholder: "Access token",
			Type:        model.SyncLoginTypeText,
			Input:       contChannel,
		},
	}

	value := <-contChannel

	events <- model.SyncEvent{Status: model.SyncStatusStarted, Lines: lines}

	if token := value.Text; token != "" {
		if ok := m.checkToken(cfg, token); !ok {
			return "", errors.New("The provided token is invalid")
		}

		m.cache.PutSecret(secretKeyAccessToken, cfg.URL, token)

		return token, nil
	}

	return "", errAbort
}

func (m *Manager) requestOAuthToken(cfg GistConfig, lines []model.SyncLine, events model.SyncEventChannel) (string, error) {
	contChannel := make(chan model.SyncInputResult)

	if token, tokenFound := m.cache.GetSecret(secretKeyAccessToken, cfg.URL); tokenFound {
		if tokenOK := m.checkToken(cfg, token); tokenOK {
			return token, nil
		} else {
			log.Info().Msgf("Stored token for %s is invalid. Delete it.", cfg.URL)
			m.cache.DeleteSecret(secretKeyAccessToken, cfg.URL)
			lines = append(lines, model.SyncLine{Type: model.SyncLineTypeError, Value: "The current token is invalid"})
		}
	}

	flow := &oauth.Flow{
		Host:      oauth.GitHubHost(cfg.hostURL()),
		ClientID:  stringutil.StringOrDefault(cfg.OAuthClientID, defaultOAuthClientID),
		Scopes:    []string{"gist"},
		BrowseURL: m.browseURL,
		DisplayCode: func(userCode string, uri string) error {
			content := uimsg.ManagerOauthDeviceFlow(cfg.hostURL(), userCode)
			events <- model.SyncEvent{
				Status: model.SyncStatusStarted,
				Lines:  lines,
				Login: &model.SyncInput{
					Content: &content,
					Type:    model.SyncLoginTypeContinue,
					Input:   contChannel,
				},
			}

			if x := <-contChannel; x.Abort {
				return errAbort
			}

			return nil
		},
	}

	accessToken, err := flow.DetectFlow()
	if err != nil {
		return "", err
	}

	events <- model.SyncEvent{Status: model.SyncStatusStarted, Lines: lines}

	if ok := m.checkToken(cfg, accessToken.Token); !ok {
		if ok = m.checkToken(cfg, accessToken.Token); !ok {
			return "", errors.New("The provided token is invalid")
		}
	}

	m.cache.PutSecret(secretKeyAccessToken, cfg.URL, accessToken.Token)

	return accessToken.Token, nil
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
