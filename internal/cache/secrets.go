package cache

import "github.com/zalando/go-keyring"

const (
	defaultSecretUser = ""
	noSecret          = ""
)

func (c *cacheImpl) GetSecret(key SecretKey) (string, bool) {
	value, err := keyring.Get(string(key), defaultSecretUser)
	switch {
	case err == keyring.ErrNotFound:
		return noSecret, false
	case err != nil:
		panic(err)
	default:
		return value, true
	}
}

func (c *cacheImpl) PutSecret(key SecretKey, secret string) {
	if err := keyring.Set(string(key), defaultSecretUser, secret); err != nil {
		panic(err)
	}
}
