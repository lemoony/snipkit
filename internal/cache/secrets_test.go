package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zalando/go-keyring"
)

const testSecretKey = SecretKey("snipkit test key")

func Test_PutAndGetSecret(t *testing.T) {
	defer func() {
		if err := keyring.Delete(string(testSecretKey), defaultSecretUser); err != nil {
			panic(err)
		}
	}()

	cache := New()

	s, ok := cache.GetSecret(testSecretKey)
	assert.False(t, ok)
	assert.Equal(t, noSecret, s)

	cache.PutSecret(testSecretKey, "password1")
	s, ok = cache.GetSecret(testSecretKey)
	assert.True(t, ok)
	assert.Equal(t, "password1", s)

	cache.PutSecret(testSecretKey, "password2")
	s, ok = cache.GetSecret(testSecretKey)
	assert.True(t, ok)
	assert.Equal(t, "password2", s)
}
