package system

import (
	"testing"

	"emperror.dev/errors"
	"github.com/stretchr/testify/assert"
)

func Test_IsErrFileSystem(t *testing.T) {
	assert.True(t, errors.Is(
		NewErrFileSystem(errors.New("root error"), "/path/to", "failed to do something"),
		ErrFileSystem{},
	))

	assert.False(t, errors.Is(errors.New("root error"), ErrFileSystem{}))
}
