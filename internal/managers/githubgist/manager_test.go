package githubgist

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/utils/testutil"
)

func Test_GetInfo(t *testing.T) {
	config := Config{
		Enabled: true,
		Gists: []GistConfig{
			{
				Enabled:              true,
				URL:                  "https://api.github.com/users/<foo-user>/gists",
				AuthenticationMethod: AuthMethodNone,
			},
		},
	}

	provider, err := NewManager(WithSystem(testutil.NewTestSystem()), WithConfig(config))
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
