package cache

type SecretKey string

type Cache interface {
	GetSecret(key SecretKey, account string) (string, bool)
	PutSecret(key SecretKey, account string, secret string)
	DeleteSecret(key SecretKey, account string)
}

type cacheImpl struct{}

func New() Cache {
	return &cacheImpl{}
}
