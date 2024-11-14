package cache

import (
	"fmt"
	"path"

	"github.com/phuslu/log"
	"github.com/spf13/afero"
	"github.com/zalando/go-keyring"
)

const noSecret = ""

func (c *cacheImpl) GetSecret(key SecretKey, account string) (string, bool) {
	if c.plainFileSecretsEnabled {
		return c.getSecretFromFile(key, account)
	} else {
		return c.getSecretFromKeyring(key, account)
	}
}

func (c *cacheImpl) getSecretFromKeyring(key SecretKey, account string) (string, bool) {
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

func (c *cacheImpl) getSecretFromFile(key SecretKey, account string) (string, bool) {
	value, err := afero.ReadFile(c.system.Fs, c.secretFilePath(key, account))
	if err != nil {
		return noSecret, false
	}
	return string(value), true
}

func (c *cacheImpl) PutSecret(key SecretKey, account, secret string) {
	if c.plainFileSecretsEnabled {
		secretPath := c.secretFilePath(key, account)
		c.system.CreatePath(secretPath)
		c.system.WriteFile(secretPath, []byte(secret))
	} else {
		if err := keyring.Set(key.service(), account, secret); err != nil {
			panic(err)
		}
	}
}

func (c *cacheImpl) DeleteSecret(key SecretKey, account string) {
	if c.plainFileSecretsEnabled {
		if err := c.system.Fs.RemoveAll(c.secretFilePath(key, account)); err != nil {
			log.Warn().Str("key", string(key)).Err(err).Msg("Failed to delete secret")
		}
	} else {
		if err := keyring.Delete(key.service(), account); err != nil && err != keyring.ErrNotFound {
			panic(err)
		}
	}
}

func (c *cacheImpl) secretFilePath(key SecretKey, account string) string {
	return path.Join(c.system.HomeDir(), ".secrets", account, string(key))
}

func (s SecretKey) service() string {
	return fmt.Sprintf("Snipkit %s", string(s))
}
