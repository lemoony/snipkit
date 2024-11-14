package githubgist

import (
	"fmt"
	"net/http"
	"os"
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
	testDataGitHubDataPath   = "testdata/github_get_user_gists_response.json"
	testDataGitHubDeviceCode = "testdata/github_device_code.json"
	testDataGitHubOAuthToken = "testdata/github_oauth_access_token.json"
	testGitHubRawURL         = "https://gist.github.test/lemoony/6e9855e7234se158b6414c/raw/52a5fd68a0fffd06297100d77c41/test-file.sh"
)

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

func Test_GetSnippets(t *testing.T) {
	testStore := expectedStoreForTestData()

	cacheMock := mocks.Cache{}
	cacheMock.On("GetData", storeKey).Return(testStore.serialize(), true)

	manager := &Manager{cache: &cacheMock, config: Config{Enabled: true, Gists: []GistConfig{
		{Enabled: true, URL: testGistURL, AuthenticationMethod: AuthMethodNone, IncludeTags: []string{"foo"}},
	}}}

	snippets := manager.GetSnippets()
	assert.Len(t, snippets, 1)
	assert.Equal(t, snippets[0].GetContent(), "foo")
}

// Scenario: Auth method is none.
// Expected: No token check is required.
func Test_Sync_noAuth(t *testing.T) {
	defer gock.Off()

	cacheMock := mocks.Cache{}
	cacheMock.On("GetData", storeKey).Return(nil, false)
	cacheMock.On("PutData", storeKey, mock.Anything).Return()

	mockGitHubGists(t)
	mockGitHubRawSnippetURL()

	manager := &Manager{cache: &cacheMock, config: Config{Enabled: true, Gists: []GistConfig{
		{Enabled: true, URL: testGistURL, AuthenticationMethod: AuthMethodNone},
	}}}

	eventChannel := make(model.SyncEventChannel)

	go func() {
		defer close(eventChannel)
		manager.Sync(eventChannel)
	}()

	for event := range eventChannel {
		t.Logf("Received event: %v\n", event)
	}

	cacheMock.AssertCalled(t, "PutData", storeKey, expectedStoreForTestData().serialize())
}

// Scenario: Auth method is to token and no token was provided previously.
// Expected: UI should prompt for token and call to GitHub is using the same token.
func Test_Sync_patAuth(t *testing.T) {
	defer gock.Off()

	cacheMock := mocks.Cache{}
	cacheMock.On("GetSecret", secretKeyAccessToken, testGistURL).Return("", false)
	cacheMock.On("GetData", storeKey).Return(nil, false)
	cacheMock.On("PutData", storeKey, mock.Anything).Return()
	cacheMock.On("PutSecret", secretKeyAccessToken, testGistURL, testToken).Return()

	mockGitHubTokenCheck(true, testToken)
	mockGitHubGists(t)
	mockGitHubRawSnippetURL()

	manager := &Manager{cache: &cacheMock, config: Config{Enabled: true, Gists: []GistConfig{
		{Enabled: true, URL: testGistURL, AuthenticationMethod: AuthMethodPAT},
	}}}

	eventChannel := make(model.SyncEventChannel)

	go func() {
		defer close(eventChannel)
		manager.Sync(eventChannel)
	}()

	for event := range eventChannel {
		t.Logf("Received event: %v\n", event)
		// test needs to provide a token for github (fake user input)
		if login := event.Login; login != nil {
			login.Input <- model.SyncInputResult{Text: testToken}
		}
	}

	cacheMock.AssertCalled(t, "PutData", storeKey, expectedStoreForTestData().serialize())
}

// Scenario: Token is expired (github returns 401 code) and user is prompted for new token but aborts.
// Expected: No cache update; token secret should be removed.
func Test_Sync_patAuth_expired_abort(t *testing.T) {
	defer gock.Off()

	const expiredToken = "expired_token"

	cacheMock := mocks.Cache{}
	cacheMock.On("GetData", storeKey).Return(nil, false)
	cacheMock.On("GetSecret", secretKeyAccessToken, testGistURL).Return(expiredToken, true)
	cacheMock.On("DeleteSecret", secretKeyAccessToken, testGistURL).Return()

	mockGitHubTokenCheck(false, expiredToken)

	manager := &Manager{cache: &cacheMock, config: Config{Enabled: true, Gists: []GistConfig{
		{Enabled: true, URL: testGistURL, AuthenticationMethod: AuthMethodPAT},
	}}}

	eventChannel := make(model.SyncEventChannel)

	go func() {
		defer close(eventChannel)
		manager.Sync(eventChannel)
	}()

	didReceiveAbort := false
	for event := range eventChannel {
		t.Logf("Received event: %v\n", event)
		// test needs to provide a token for github (fake user input)
		if login := event.Login; login != nil {
			login.Input <- model.SyncInputResult{Abort: true}
		}

		if event.Status == model.SyncStatusAborted {
			didReceiveAbort = true
		}
	}

	assert.True(t, didReceiveAbort)

	cacheMock.AssertCalled(t, "DeleteSecret", secretKeyAccessToken, testGistURL)
	cacheMock.AssertNotCalled(t, "PutData", storeKey, mock.Anything)
}

// Scenario: Auth method is to oauth device flow and no token was provided previously.
// Expected: UI should prompt for pressing enter to continue and GitHub device flow is triggered.
func Test_Sync_oauthDeviceFlow(t *testing.T) {
	defer gock.Off()

	cacheMock := mocks.Cache{}
	cacheMock.On("GetSecret", secretKeyAccessToken, testGistURL).Return("", false)
	cacheMock.On("PutSecret", secretKeyAccessToken, mock.Anything, mock.Anything).Return()
	cacheMock.On("GetData", storeKey).Return(nil, false)
	cacheMock.On("PutData", storeKey, mock.Anything).Return()

	mockGitHubDeviceCode(t)
	mockGitHubBearerToken(t)
	mockGitHubTokenCheck(true, testToken)
	mockGitHubGists(t)
	mockGitHubRawSnippetURL()

	didCallBrowseURL := false
	browseURLFunc := func(s string) error {
		didCallBrowseURL = true
		return nil
	}

	manager := &Manager{cache: &cacheMock, config: Config{Enabled: true, Gists: []GistConfig{
		{Enabled: true, URL: testGistURL, AuthenticationMethod: AuthMethodOAuthDeviceFlow},
	}}, browseURL: browseURLFunc}

	eventChannel := make(model.SyncEventChannel)

	go func() {
		defer close(eventChannel)
		manager.Sync(eventChannel)
	}()

	for event := range eventChannel {
		t.Logf("Received event: %v\n", event)
		// test needs to provide a token for github (fake user input)
		if login := event.Login; login != nil {
			login.Input <- model.SyncInputResult{Text: testToken}
		}
	}

	assert.True(t, didCallBrowseURL)

	cacheMock.AssertCalled(t, "PutSecret", secretKeyAccessToken, testGistURL, testToken)
	cacheMock.AssertCalled(t, "PutData", storeKey, expectedStoreForTestData().serialize())
}

// Scenario: Token is expired (github returns 401 code) and user is prompted for device flow but aborts
// Expected: No cache update; token secret should be removed.
func Test_Sync_oauthDeviceFlow_abort(t *testing.T) {
	defer gock.Off()

	const expiredToken = "toke_expired"

	cacheMock := mocks.Cache{}
	cacheMock.On("GetSecret", secretKeyAccessToken, testGistURL).Return(expiredToken, true)
	cacheMock.On("DeleteSecret", secretKeyAccessToken, testGistURL).Return()
	cacheMock.On("GetData", storeKey).Return(nil, false)

	mockGitHubTokenCheck(false, expiredToken)
	mockGitHubDeviceCode(t)
	mockGitHubBearerToken(t)
	mockGitHubGists(t)
	mockGitHubRawSnippetURL()

	didCallBrowseURL := false
	browseURLFunc := func(s string) error {
		didCallBrowseURL = true
		return nil
	}

	manager := &Manager{cache: &cacheMock, config: Config{Enabled: true, Gists: []GistConfig{
		{Enabled: true, URL: testGistURL, AuthenticationMethod: AuthMethodOAuthDeviceFlow},
	}}, browseURL: browseURLFunc}

	eventChannel := make(model.SyncEventChannel)

	go func() {
		defer close(eventChannel)
		manager.Sync(eventChannel)
	}()

	for event := range eventChannel {
		t.Logf("Received event: %v\n", event)
		if login := event.Login; login != nil {
			login.Input <- model.SyncInputResult{Abort: true}
		}
	}

	assert.False(t, didCallBrowseURL)

	cacheMock.AssertCalled(t, "DeleteSecret", secretKeyAccessToken, testGistURL)
	cacheMock.AssertNotCalled(t, "PutData")
}

// Scenario: Sync is triggerd and the cache already contains entries. ETag from github does not change (status 304)
// Expected: The same data is put into the store as it was retrieved previously.
func Test_Sync_ifNoneMatch(t *testing.T) {
	defer gock.Off()

	cachedStore := expectedStoreForTestData()

	cacheMock := mocks.Cache{}
	cacheMock.On("GetData", storeKey).Return(cachedStore.serialize(), true)
	cacheMock.On("PutData", storeKey, mock.Anything).Return()

	gock.New(fmt.Sprintf("https://api.%s", testHost)).
		MatchHeader("If-None-Match", cachedStore.Gists[0].ETag).
		Get(fmt.Sprintf("users/%s/gists", testUser)).
		Reply(http.StatusNotModified).
		SetHeader("etag", cachedStore.Gists[0].ETag)

	manager := &Manager{cache: &cacheMock, config: Config{Enabled: true, Gists: []GistConfig{
		{Enabled: true, URL: testGistURL, AuthenticationMethod: AuthMethodNone},
	}}}

	eventChannel := make(model.SyncEventChannel)
	go func() {
		defer close(eventChannel)
		manager.Sync(eventChannel)
	}()

	for event := range eventChannel {
		t.Logf("Received event: %v\n", event)
	}

	cacheMock.AssertCalled(t, "PutData", storeKey, cachedStore.serialize())
}

// Scenario: A single gist file has changed, so the etag values are different than the ones previously stored.
// Expected: Update the cache with new etag values.
func Test_Sync_ifNoneMatch_forSingleFile(t *testing.T) {
	defer gock.Off()

	const updatedGistEtag = "etag_updated"
	const updatedFileEtag = "etag_file_updated"
	cachedStore := expectedStoreForTestData()

	cacheMock := mocks.Cache{}
	cacheMock.On("GetData", storeKey).Return(cachedStore.serialize(), true)
	cacheMock.On("PutData", storeKey, mock.Anything).Return()

	gock.New(fmt.Sprintf("https://api.%s", testHost)).
		MatchHeader("If-None-Match", cachedStore.Gists[0].ETag).
		Get(fmt.Sprintf("users/%s/gists", testUser)).
		Reply(http.StatusOK).
		SetHeader("etag", updatedGistEtag).
		JSON(readTestdata(t, testDataGitHubDataPath))

	gock.New(testGitHubRawURL).
		Get("").
		MatchHeader("If-None-Match", cachedStore.Gists[0].RawSnippets[0].ETag).
		Reply(http.StatusOK).
		SetHeader("etag", updatedFileEtag).
		BodyString("foo")

	manager := &Manager{cache: &cacheMock, config: Config{Enabled: true, Gists: []GistConfig{
		{Enabled: true, URL: testGistURL, AuthenticationMethod: AuthMethodNone},
	}}}

	eventChannel := make(model.SyncEventChannel)
	go func() {
		defer close(eventChannel)
		manager.Sync(eventChannel)
	}()

	for event := range eventChannel {
		t.Logf("Received event: %v\n", event)
	}

	updatedStore := *cachedStore
	updatedStore.Gists[0].ETag = updatedGistEtag
	updatedStore.Gists[0].RawSnippets[0].ETag = updatedFileEtag
	updatedStore.Gists[0].RawSnippets[0].Content = []byte("foo")

	cacheMock.AssertCalled(t, "PutData", storeKey, updatedStore.serialize())
}

func Test_SaveAssistantSnippet(t *testing.T) {
	assert.PanicsWithError(t, "Not implemented", func() {
		Manager{}.SaveAssistantSnippet("", "foo.sh", []byte("dummy content"))
	})
}

func readTestdata(t *testing.T, path string) string {
	t.Helper()
	contents, err := os.ReadFile(path)
	assert.NoError(t, err)
	return string(contents)
}

func expectedStoreForTestData() *store {
	return &store{
		Version: storeVersion,
		Gists: []gistStore{
			{
				URL:  testGistURL,
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

func mockGitHubTokenCheck(valid bool, token string) {
	status := http.StatusOK
	if !valid {
		status = http.StatusUnauthorized
	}

	gock.New(testAPIURL).
		MatchHeader("Authorization", fmt.Sprintf("token %s", token)).
		Head(fmt.Sprintf("users/%s/gists", testUser)).
		Reply(status)
}

func mockGitHubGists(t *testing.T) {
	t.Helper()
	gock.New(testAPIURL).
		Get(fmt.Sprintf("users/%s/gists", testUser)).
		Reply(200).
		SetHeader("etag", "test_etag_value").
		JSON(readTestdata(t, testDataGitHubDataPath))
}

func mockGitHubRawSnippetURL() {
	gock.New(testGitHubRawURL).Get("").Reply(200).
		SetHeader("etag", "raw_snippet_etag_value").
		BodyString("foo")
}

func mockGitHubDeviceCode(t *testing.T) {
	t.Helper()
	gock.New(fmt.Sprintf("https://%s", testHost)).
		BodyString(fmt.Sprintf(`client_id=%s&scope=gist`, defaultOAuthClientID)).
		Post("login/device/code").
		Reply(200).
		JSON(readTestdata(t, testDataGitHubDeviceCode))
}

func mockGitHubBearerToken(t *testing.T) {
	t.Helper()
	gock.New(fmt.Sprintf("https://%s", testHost)).
		Post("login/oauth/access_token").
		Reply(200).
		JSON(fmt.Sprintf(readTestdata(t, testDataGitHubOAuthToken), testToken))
}
