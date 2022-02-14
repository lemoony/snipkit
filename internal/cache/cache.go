package cache

import "github.com/lemoony/snipkit/internal/utils/system"

type (
	SecretKey string
	DataKey   string
)

type Cache interface {
	GetSecret(key SecretKey, account string) (string, bool)
	PutSecret(key SecretKey, account string, secret string)
	DeleteSecret(key SecretKey, account string)
	PutData(key DataKey, data []byte)
	GetData(key DataKey) ([]byte, bool)
}

type cacheImpl struct {
	system *system.System
}

func New(s *system.System) Cache {
	return &cacheImpl{
		system: s,
	}
}
