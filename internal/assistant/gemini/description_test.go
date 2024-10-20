package gemini

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescription(t *testing.T) {
	description := Description(&Config{
		Enabled: true,
	})

	assert.Equal(t, Key, description.Key)
	assert.Equal(t, true, description.Enabled)
}
