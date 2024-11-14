package cache

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/utils/testutil"
)

const (
	testSecretKey = SecretKey("snipkit test key")
	testAccount1  = "github.com/testuser"
	testAccount2  = "github.com/otheruser"
)

//nolint:funlen // test function for yaml config is allowed to be too long
func Test_PutAndGetSecret(t *testing.T) {
	tests := []struct {
		name                    string
		plainFileSecretsEnabled bool
		skip                    bool
	}{
		{"KeyringSecrets", false, runtime.GOOS != "darwin"},
		{"PlainFileSecrets", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				// Setting up gnome-keyring in a headless environment is not easy...
				t.Skip("Test is skipped for the present OS")
			}

			cache := &cacheImpl{system: testutil.NewTestSystem()}

			if tt.plainFileSecretsEnabled {
				cache.EnablePlainFileSecrets()
			}

			defer func() {
				cache.DeleteSecret(testSecretKey, testAccount1)
				cache.DeleteSecret(testSecretKey, testAccount2)
			}()

			cache.DeleteSecret(testSecretKey, testAccount1)
			cache.DeleteSecret(testSecretKey, testAccount2)

			time.Sleep(time.Millisecond * 50)

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
		})
	}
}
