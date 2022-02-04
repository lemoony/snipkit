package cache

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

const noSecret = ""

func (c *cacheImpl) GetSecret(key SecretKey, account string) (string, bool) {
	value, err := keyring.Get(key.service(), account)
	switch {
	case err == keyring.ErrNotFound:
		return noSecret, false
	case err != nil:
		panic(err)
	default:
		return value, true
	}
}

func (c *cacheImpl) PutSecret(key SecretKey, account, secret string) {
	if err := keyring.Set(key.service(), account, secret); err != nil {
		panic(err)
	}
}

func (c *cacheImpl) DeleteSecret(key SecretKey, account string) {
	if err := keyring.Delete(key.service(), account); err != nil && err != keyring.ErrNotFound {
		panic(err)
	}
}

func (s SecretKey) service() string {
	return fmt.Sprintf("Snipkit %s", string(s))
}
