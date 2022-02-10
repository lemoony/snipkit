package githubgist

import (
	"encoding/json"

	"github.com/phuslu/log"

	"github.com/lemoony/snipkit/internal/cache"
)

const storeKey = cache.DataKey("github_gist_cache")

type store struct {
	Version string      `json:"version"`
	Gists   []gistStore `json:"gists"`
}

type gistStore struct {
	URL         string       `json:"url"`
	ETag        string       `json:"ETag"`
	RawSnippets []rawSnippet `json:"RawSnippets"`
}

type rawSnippet struct {
	ID          string `json:"id"`
	Filename    string `json:"filename"`
	Content     []byte `json:"content"`
	ETag        string `json:"etag"`
	Pubic       bool   `json:"public"`
	Description string `json:"description"`
}

func (m *Manager) loadFromCache() (map[string]*gistStore, bool) {
	if bytes, ok := m.cache.GetData(storeKey); !ok {
		return nil, false
	} else {
		var s store
		if err := json.Unmarshal(bytes, &s); err != nil {
			panic(err)
		}

		result := map[string]*gistStore{}
		for i := range s.Gists {
			result[s.Gists[i].URL] = &s.Gists[i]
		}
		return result, true
	}
}

func (m *Manager) getStoreFromCache() *store {
	raw, ok := m.cache.GetData(storeKey)
	if !ok {
		return nil
	}

	var resultStore store
	if !resultStore.deserialize(raw) {
		return nil
	}

	return &resultStore
}

func (m *Manager) storeInCache(all map[string]*gistStore) {
	var s store
	s.Version = "1.0.0"
	s.Gists = make([]gistStore, len(all))
	i := 0
	for k := range all {
		s.Gists[i] = *all[k]
	}

	bytes := s.serialize()
	m.cache.PutData(storeKey, bytes)
}

func (c *store) serialize() []byte {
	if bytes, err := json.Marshal(c); err != nil {
		panic(err)
	} else {
		return bytes
	}
}

func (c *store) deserialize(bytes []byte) bool {
	if err := json.Unmarshal(bytes, c); err != nil {
		log.Info().Err(err).Msg("store invalid")
		return false
	}
	return true
}

func (c *store) getGists(cfg GistConfig) *gistStore {
	for i := range c.Gists {
		if c.Gists[i].URL == cfg.URL {
			return &c.Gists[i]
		}
	}
	return nil
}
