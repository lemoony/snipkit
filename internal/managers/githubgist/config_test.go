package githubgist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AutoDiscoveryConfig(t *testing.T) {
	cfg := AutoDiscoveryConfig()
	assert.NotNil(t, cfg)
	assert.False(t, cfg.Enabled)
	assert.Len(t, cfg.Gists, 1)

	gistConfig := cfg.Gists[0]
	assert.Equal(t, "gist.github.com/<USERNAME>", gistConfig.URL)
	assert.Equal(t, "https://api.github.com/users/<USERNAME>/gists", gistConfig.apiURL())

	assert.Equal(t, &gistConfig, cfg.getGistConfig("gist.github.com/<USERNAME>"))
}

func Test_getGistConfig_unknown(t *testing.T) {
	assert.Nil(t, AutoDiscoveryConfig().getGistConfig("gist.github.com/foouser"))
}

func Test_config_getAPIUrl(t *testing.T) {
	gistConfig := GistConfig{URL: "gist.github.com/<USERNAME>"}
	gistConfig.URL = "http://someurl.com"
	assert.Panics(t, func() {
		gistConfig.apiURL()
	})
}

func Test_config_invalidURL(t *testing.T) {
	gistConfig := GistConfig{URL: "someurl.com"}
	gistConfig.URL = "http://someurl.com"
	assert.Panics(t, func() {
		gistConfig.apiURL()
	})
}
