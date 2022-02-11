package githubgist

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/h2non/gock.v1"

	"github.com/lemoony/snipkit/internal/cache"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/testutil"
	mocks "github.com/lemoony/snipkit/mocks/cache"
)

const (
	testGitHubTestDataPath = "testdata/github_get_user_gists_response.json"
	testGitHubRawURL       = "https://gist.github.test/lemoony/6e9855e7234se158b6414c/raw/52a5fd68a0fffd06297100d77c41/test-file.sh"
)

var testURL = fmt.Sprintf("gist.%s/%s", testHost, testUser)

func Test_GetInfo(t *testing.T) {
	config := Config{
		Enabled: true,
		Gists: []GistConfig{
			{
				Enabled:              true,
				URL:                  "github.com/foo-user",
				AuthenticationMethod: AuthMethodNone,
			},
		},
	}

	system := testutil.NewTestSystem()
	provider, err := NewManager(WithSystem(system), WithConfig(config), WithCache(cache.New(system)))
	assert.NoError(t, err)

	info := provider.Info()

	assert.Len(t, info, 3)

	assert.Equal(t, info[0].Key, "GitHub Gist enabled")
	assert.Equal(t, info[0].Value, "true")
	assert.False(t, info[0].IsError)

	assert.Equal(t, info[1].Key, "GitHub Gist number of URLs")
	assert.Equal(t, info[1].Value, "1")
	assert.False(t, info[1].IsError)

	assert.Equal(t, info[2].Key, "GitHub Gist total number of snippets")
	assert.Equal(t, info[2].Value, "0")
	assert.False(t, info[2].IsError)
}

func Test_Key(t *testing.T) {
	assert.Equal(t, Key, Manager{}.Key())
}

func Test_NewManager_disabled(t *testing.T) {
	manager, err := NewManager(WithConfig(Config{Enabled: false}))
	assert.NoError(t, err)
	assert.Nil(t, manager)
}

func Test_Sync_noAuth(t *testing.T) {
	defer gock.Off()

	cacheMock := mocks.Cache{}
	cacheMock.On("GetData", storeKey).Return(nil, false)
	cacheMock.On("PutData", storeKey, mock.Anything).Return()

	gock.New(fmt.Sprintf("https://api.%s", testHost)).
		Get(fmt.Sprintf("users/%s/gists", testUser)).
		Reply(200).
		SetHeader("etag", "test_etag_value").
		JSON(readTestdata(t, testGitHubTestDataPath))

	gock.New(testGitHubRawURL).Get("").Reply(200).
		SetHeader("etag", "raw_snippet_etag_value").
		BodyString("foo")

	manager := &Manager{cache: &cacheMock, config: Config{Enabled: true, Gists: []GistConfig{
		{Enabled: true, URL: testURL, AuthenticationMethod: AuthMethodNone},
	}}}

	eventChannel := make(model.SyncEventChannel)

	doneSync := make(chan struct{})
	go func() {
		defer close(doneSync)
		success := manager.Sync(eventChannel)
		assert.True(t, success)
	}()

	for event := range eventChannel {
		t.Logf("Received event: %v\n", event)
	}

	<-doneSync

	cacheMock.AssertCalled(t, "PutData", storeKey, expectedStoreForTestData().serialize())
}

func readTestdata(t *testing.T, path string) string {
	t.Helper()
	contents, err := ioutil.ReadFile(path)
	assert.NoError(t, err)
	return string(contents)
}

func expectedStoreForTestData() *store {
	return &store{
		Version: storeVersion,
		Gists: []gistStore{
			{
				URL:  testURL,
				ETag: "test_etag_value",
				RawSnippets: []rawSnippet{
					{
						ID:          "testsnippetid-test-file.sh",
						ETag:        "raw_snippet_etag_value",
						Description: "Echo Something #foo",
						Filename:    "test-file.sh",
						FilesInGist: 1,
						Pubic:       true,
						Language:    "Shell",
						Content:     []byte("foo"),
					},
				},
			},
		},
	}
}
