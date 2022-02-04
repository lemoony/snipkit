//go:build darwin
// +build darwin

package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testSecretKey = SecretKey("snipkit test key")
	testAccount1  = "github.com/testuser"
	testAccount2  = "github.com/otheruser"
)

// For now this test is only executed for MacOS since setting up gnome-keyring in a headless environment is not that easy...
func Test_PutAndGetSecret(t *testing.T) {
	cache := New()

	defer func() {
		cache.DeleteSecret(testSecretKey, testAccount1)
		cache.DeleteSecret(testSecretKey, testAccount2)
	}()

	s, ok := cache.GetSecret(testSecretKey, testAccount1)
	assert.False(t, ok)
	assert.Equal(t, noSecret, s)

	s, ok = cache.GetSecret(testSecretKey, testAccount2)
	assert.False(t, ok)
	assert.Equal(t, noSecret, s)

	cache.PutSecret(testSecretKey, testAccount1, "password1")
	cache.PutSecret(testSecretKey, testAccount2, "foo1")
	s, ok = cache.GetSecret(testSecretKey, testAccount1)
	assert.True(t, ok)
	assert.Equal(t, "password1", s)

	s, ok = cache.GetSecret(testSecretKey, testAccount2)
	assert.True(t, ok)
	assert.Equal(t, "foo1", s)

	cache.PutSecret(testSecretKey, testAccount1, "password2")
	cache.PutSecret(testSecretKey, testAccount2, "foo2")

	s, ok = cache.GetSecret(testSecretKey, testAccount1)
	assert.True(t, ok)
	assert.Equal(t, "password2", s)

	s, ok = cache.GetSecret(testSecretKey, testAccount2)
	assert.True(t, ok)
	assert.Equal(t, "foo2", s)

	cache.DeleteSecret(testSecretKey, testAccount1)
	if _, exists := cache.GetSecret(testSecretKey, testAccount1); exists {
		assert.False(t, exists)
	}
	if _, exists := cache.GetSecret(testSecretKey, testAccount2); !exists {
		assert.True(t, exists)
	}

	cache.DeleteSecret(testSecretKey, testAccount2)
	if _, exists := cache.GetSecret(testSecretKey, testAccount2); exists {
		assert.False(t, exists)
	}
}
