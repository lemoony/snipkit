package githubgist

import (
	"fmt"
	"net/http"
	"testing"

	"emperror.dev/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"

	"github.com/lemoony/snipkit/internal/utils/assertutil"
)

const (
	testToken = "test_token"
	testHost  = "github.test"
	testUser  = "foouser"

	testNoETag         = ""
	testNoResponseEtag = ""
)

var (
	testAPIURL  = fmt.Sprintf("https://api.%s", testHost)
	testGistURL = fmt.Sprintf("gist.%s/%s", testHost, testUser)
)

func Test_checkToken_valid(t *testing.T) {
	testCheckToken(t, 200, nil, true)
}

func Test_checkToken_invalid(t *testing.T) {
	testCheckToken(t, 401, nil, false)
}

func Test_checkToken_invalidStatusCode(t *testing.T) {
	_ = assertutil.AssertPanicsWithError(t, errUnexpected, func() {
		testCheckToken(t, 400, nil, false)
	})
}

func Test_checkToken_requestError(t *testing.T) {
	assert.Panics(t, func() {
		testCheckToken(t, 400, errors.New("foo error"), false)
	})
}

func testCheckToken(t *testing.T, statusCode int, err error, expectedValid bool) {
	t.Helper()
	defer gock.Off()

	cfg := Config{
		Enabled: true,
		Gists: []GistConfig{
			{
				URL:                  fmt.Sprintf("gist.%s/%s", testHost, testUser),
				AuthenticationMethod: AuthMethodPAT,
			},
		},
	}

	m, _ := NewManager(WithConfig(cfg))

	request := gock.New(fmt.Sprintf("https://api.%s", testHost)).
		MatchHeader("Authorization", fmt.Sprintf("token %s", testToken)).
		Head(fmt.Sprintf("users/%s/gists", testUser))

	if err != nil {
		request.ReplyError(err)
	} else {
		request.Reply(statusCode)
	}

	assert.Equal(t, expectedValid, m.checkToken(cfg.Gists[0], testToken))
}

func Test_getRawGist(t *testing.T) {
	defer gock.Off()

	const path = "some/raw/gist-id"

	manager := prepareGetRawResponse(path, "", 200, testNoResponseEtag, `foo: test`)

	actualResponse := manager.getRawGist(
		fmt.Sprintf("https://api.%s/%s", testHost, path), "etag-value", testToken,
	)
	assert.True(t, actualResponse.hasUpdates)
	assert.NotNil(t, actualResponse.rawContent)
	assert.Nil(t, actualResponse.gistsResponse)
	assert.Equal(t, []byte(`foo: test`), *actualResponse.rawContent)
}

func Test_getRawGist_NotModified(t *testing.T) {
	defer gock.Off()

	const path = "some/raw/gist-id"
	const etagValue = "etag-value"

	manager := prepareGetRawResponse(path, etagValue, http.StatusNotModified, testNoResponseEtag, `foo: test`)

	actualResponse := manager.getRawGist(
		fmt.Sprintf("https://api.%s/%s", testHost, path), etagValue, testToken,
	)
	assert.False(t, actualResponse.hasUpdates)
	assert.Nil(t, actualResponse.rawContent)
	assert.Nil(t, actualResponse.gistsResponse)
}

func Test_getRawGist_WeakEtag(t *testing.T) {
	defer gock.Off()

	const path = "some/raw/gist-id"

	manager := prepareGetRawResponse(path, "", http.StatusOK, `W/"weaketag"`, "foo")

	response := manager.getRawGist(
		fmt.Sprintf("https://api.%s/%s", testHost, path), "", testToken,
	)

	assert.Equal(t, "weaketag", response.etag)
}

func Test_getRawGist_BadCredentials(t *testing.T) {
	defer gock.Off()

	const path = "some/raw/gist-id"
	const response = `{
  "message": "Bad credentials",
  "documentation_url": "https://docs.github.com/rest"
}`

	manager := prepareGetRawResponse(path, "", http.StatusUnauthorized, testNoResponseEtag, response)

	err := assertutil.AssertPanicsWithError(t, errAuth, func() {
		manager.getRawGist(
			fmt.Sprintf("https://api.%s/%s", testHost, path), "", testToken,
		)
	})

	assert.Contains(t, err.Error(), response)
}

func Test_getRawGist_NoContent(t *testing.T) {
	defer gock.Off()

	const path = "some/raw/gist-id"

	manager := prepareGetRawResponse(path, "", http.StatusNoContent, testNoResponseEtag, "")

	response := manager.getRawGist(
		fmt.Sprintf("https://api.%s/%s", testHost, path), "", testToken,
	)

	assert.NotNil(t, response.rawContent)
	assert.Empty(t, response.rawContent)
}

func Test_getRawGist_Error(t *testing.T) {
	defer gock.Off()

	const path = "some/raw/gist-id"

	cfg := Config{
		Enabled: true,
		Gists: []GistConfig{
			{
				URL:                  fmt.Sprintf("gist.%s/%s", testHost, testUser),
				AuthenticationMethod: AuthMethodPAT,
			},
		},
	}
	manager, _ := NewManager(WithConfig(cfg))

	gock.New(fmt.Sprintf("https://api.%s", testHost)).
		Get(path).
		Response.SetError(errors.New("foo error"))

	assert.Panics(t, func() {
		manager.getRawGist(
			fmt.Sprintf("https://api.%s/%s", testHost, path), "", testToken,
		)
	})
}

func Test_getGists(t *testing.T) {
	manager := prepareGetRawResponse(fmt.Sprintf("users/%s/gists", testUser), testNoETag, http.StatusOK, testNoResponseEtag, `[
  {
    "files": {
      "echo.sh": {
        "filename": "private-echo-2.sh",
        "language": "Shell",
        "raw_url": "https://gist.githubusercontent.com/lemoony/d64cd2aec9517d393f130fa5f96be68c/raw/7a9d9714f668c03a733b08f9ba9a990a33520030/private-echo-2.sh"
      }
    },
    "public": false
  }
]`)

	response := manager.getGists(manager.config.Gists[0], testNoETag, testToken)

	assert.True(t, response.hasUpdates)
	assert.Nil(t, response.rawContent)
	assert.NotNil(t, response.gistsResponse, 1)
	assert.Len(t, *response.gistsResponse, 1)
}

func Test_getGists_NoUpdates(t *testing.T) {
	manager := prepareGetRawResponse(fmt.Sprintf("users/%s/gists", testUser), testNoETag, http.StatusNotModified, testNoResponseEtag, "")

	response := manager.getGists(manager.config.Gists[0], testNoETag, testToken)
	assert.False(t, response.hasUpdates)
	assert.Nil(t, response.rawContent)
	assert.Nil(t, response.gistsResponse, 1)
}

func prepareGetRawResponse(path, etag string, statusCode int, responseEtag, responseContent string) *Manager {
	cfg := Config{
		Enabled: true,
		Gists: []GistConfig{
			{
				URL:                  fmt.Sprintf("gist.%s/%s", testHost, testUser),
				AuthenticationMethod: AuthMethodPAT,
			},
		},
	}

	m, _ := NewManager(WithConfig(cfg))

	request := gock.New(fmt.Sprintf("https://api.%s", testHost)).
		MatchHeader("Authorization", fmt.Sprintf("token %s", testToken))

	if etag != "" {
		request.MatchHeader("If-None-Match", fmt.Sprintf(`"%s"`, etag))
	}

	response := request.Get(path).
		Reply(statusCode).
		JSON(responseContent)

	if responseEtag != "" {
		response.SetHeader("etag", responseEtag)
	}

	return m
}
