package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/utils/testutil"
)

func Test_PutAndGetData(t *testing.T) {
	cache := New(testutil.NewTestSystem())

	const testKey = DataKey("test_key_1")

	data, ok := cache.GetData(testKey)
	assert.False(t, ok)
	assert.Nil(t, data)

	cache.PutData(testKey, []byte("foo"))

	data, ok = cache.GetData(testKey)
	assert.True(t, ok)
	assert.Equal(t, []byte("foo"), data)
}
