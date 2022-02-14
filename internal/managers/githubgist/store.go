package githubgist

import (
	"encoding/json"

	"github.com/phuslu/log"

	"github.com/lemoony/snipkit/internal/cache"
)

const (
	storeKey     = cache.DataKey("github_gist_cache")
	storeVersion = "1.0"
)

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
	Language    string `json:"language"`
	FilesInGist int    `json:"filesInGist"`
}

func (m *Manager) getStoreFromCache() *store {
	result := &store{}
	if raw, ok := m.cache.GetData(storeKey); ok {
		result.deserialize(raw)
	}
	return result
}

func (m *Manager) storeInCache(s *store) {
	m.cache.PutData(storeKey, s.serialize())
}

func (c *store) serialize() []byte {
	if bytes, err := json.Marshal(c); err != nil {
		panic(err)
	} else {
		return bytes
	}
}

func (c *store) deserialize(bytes []byte) {
	if err := json.Unmarshal(bytes, c); err != nil {
		log.Warn().Err(err).Msg("store invalid")
	}
}

func (c *store) getGists(cfg GistConfig) *gistStore {
	for i := range c.Gists {
		if c.Gists[i].URL == cfg.URL {
			return &c.Gists[i]
		}
	}
	return nil
}
