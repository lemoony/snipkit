package githubgist

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mocks "github.com/lemoony/snipkit/mocks/cache"
)

func Test_getStoreFromCache(t *testing.T) {
	const cachedStore = `{
    "version": "1.0.0",
    "gists": [
        {
            "url": "gist.github.com/lemoony",
            "ETag": "591c38f5247596e7f6d4f3093d45586127b57f1e307f8dab9477aa10a3cfd32a",
            "RawSnippets": [
                {
                    "id": "d64cd2aec9517d393f130fa5f96be68c-private-echo-2.sh",
                    "filename": "private-echo-2.sh",
                    "content": "ZWNobyAiYW5vdGhlciBwcml2YXRlIGdpc3QiCiNibGE=",
                    "etag": "456e3b32afefc84ddb207d58e5cda59102cb5c0a99882a7e5775d1ddeac7ac72",
                    "public": false,
                    "description": "Some private snippets #1"
                }
            ]
        }
    ]
}`

	cache := mocks.Cache{}

	cache.On("GetData", storeKey).Return([]byte(cachedStore), true)
	manager := &Manager{cache: &cache}

	actualStore := manager.getStoreFromCache()
	assert.NotNil(t, actualStore)
	assert.Len(t, actualStore.Gists, 1)
	assert.Equal(t, actualStore.Gists[0].URL, "gist.github.com/lemoony")
	assert.Len(t, actualStore.Gists[0].RawSnippets, 1)
}

func Test_loadFromCache_Empty(t *testing.T) {
	cache := mocks.Cache{}
	cache.On("GetData", storeKey).Return(nil, false)
	manager := &Manager{cache: &cache}
	assert.Empty(t, manager.getStoreFromCache().Gists)
}

func Test_loadFromCache_InvalidContent(t *testing.T) {
	tests := []struct {
		name                     string
		content                  string
		expectedStoreGistsLength int
	}{
		{name: "no content", content: ``, expectedStoreGistsLength: 0},
		{name: "wrong json schema", content: `{"foo": "invalid"}`, expectedStoreGistsLength: 0},
		{name: "invalid json", content: `foo`, expectedStoreGistsLength: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := mocks.Cache{}
			cache.On("GetData", storeKey).Return([]byte(tt.content), true)
			manager := &Manager{cache: &cache}
			cachedStore := manager.getStoreFromCache()
			assert.Len(t, cachedStore.Gists, tt.expectedStoreGistsLength)
		})
	}
}

func Test_storeInCache(t *testing.T) {
	cache := mocks.Cache{}
	cache.On("PutData", storeKey, mock.Anything).Return()
	manager := &Manager{cache: &cache}

	s := &store{Version: storeVersion, Gists: []gistStore{
		{URL: "test.url", ETag: "etag1", RawSnippets: []rawSnippet{{ID: "id", ETag: "etag2"}}},
	}}
	manager.storeInCache(s)

	cache.AssertNumberOfCalls(t, "PutData", 1)
	cache.AssertCalled(t, "PutData", storeKey, []byte(
		`{"version":"1.0","gists":[{"url":"test.url","ETag":"etag1","RawSnippets":[{"id":"id","filename":"","content":null,"etag":"etag2","public":false,"description":"","language":"","filesInGist":0}]}]}`,
	))
}

func Test_getStoreGist(t *testing.T) {
	gistStore1 := gistStore{URL: "gist.github1.com/foo"}
	gistStore2 := gistStore{URL: "gist.github2.com/foo"}
	s := &store{Version: storeVersion, Gists: []gistStore{gistStore1, gistStore2}}

	tests := []struct {
		url      string
		expected *gistStore
	}{
		{url: "gist.github1.com/foo", expected: &gistStore1},
		{url: "gist.github2.com/foo", expected: &gistStore2},
		{url: "gist.github3.com/foo", expected: nil},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			assert.Equal(t, tt.expected, s.getGists(GistConfig{URL: tt.url}))
		})
	}
}
