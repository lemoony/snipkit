package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/ui/style"
	"github.com/lemoony/snipkit/internal/ui/uimsg"
)

func Test_ErrConfigNotFound_string(t *testing.T) {
	assert.Equal(t, uimsg.ConfigNotFound("path").RenderWith(style.NoopStyle), ErrConfigNotFound{"path"}.Error())
}

func Test_ErrConfigNotFound_Is(t *testing.T) {
	assert.True(t, errors.Is(ErrConfigNotFound{"foo"}, ErrConfigNotFound{}))
	assert.False(t, errors.Is(errors.New("foo error"), ErrConfigNotFound{}))
}
