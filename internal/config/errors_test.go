package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ErrConfigNotFound_string(t *testing.T) {
	assert.Equal(t, "config not found: path", ErrConfigNotFound{"path"}.Error())
}

func Test_ErrConfigNotFound_Is(t *testing.T) {
	assert.True(t, errors.Is(ErrConfigNotFound{"foo"}, ErrConfigNotFound{}))
	assert.False(t, errors.Is(errors.New("foo error"), ErrConfigNotFound{}))
}