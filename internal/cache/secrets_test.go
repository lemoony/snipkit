package cache

import (
	"testing"

	"github.com/keybase/go-keychain"
	"github.com/stretchr/testify/assert"
)

const testLabel = "snipkit_test"

func Test_PutAndGetSecret(t *testing.T) {
	defer func() {
		item := keychain.NewItem()
		item.SetSecClass(keychain.SecClassGenericPassword)
		item.SetLabel(testLabel)
		if err := keychain.DeleteItem(item); err != nil {
			panic(err)
		}
	}()

	cache := New()

	if impl, ok := cache.(*cacheImpl); !ok {
		assert.Fail(t, "cache is not of type cacheImp")
	} else {
		impl.label = testLabel
	}

	s, ok := cache.GetSecret(ServicePrivateAccessToken)
	assert.False(t, ok)
	assert.Nil(t, s)

	cache.PutSecret(ServicePrivateAccessToken, []byte("password1"))
	s, ok = cache.GetSecret(ServicePrivateAccessToken)
	assert.True(t, ok)
	assert.Equal(t, []byte("password1"), s)

	cache.PutSecret(ServicePrivateAccessToken, []byte("password2"))
	s, ok = cache.GetSecret(ServicePrivateAccessToken)
	assert.True(t, ok)
	assert.Equal(t, []byte("password2"), s)
}
