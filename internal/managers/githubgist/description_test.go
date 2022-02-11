package githubgist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Description(t *testing.T) {
	cfg := &Config{Enabled: true}

	description := Description(cfg)
	assert.NotNil(t, description)
	assert.Equal(t, "Github Gist", description.Name)
	assert.True(t, description.Enabled)
}
