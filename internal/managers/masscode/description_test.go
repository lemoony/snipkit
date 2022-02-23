package masscode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Description(t *testing.T) {
	cfg := &Config{Enabled: true}

	description := Description(cfg)
	assert.NotNil(t, description)
	assert.Equal(t, "massCode", description.Name)
	assert.True(t, description.Enabled)
}
