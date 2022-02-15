package pet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Description(t *testing.T) {
	cfg := &Config{Enabled: true}

	description := Description(cfg)
	assert.NotNil(t, description)
	assert.Equal(t, "pet - CLI Snippet Manager", description.Name)
	assert.True(t, description.Enabled)
}
