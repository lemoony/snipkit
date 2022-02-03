package cache

type SecretKey string

type Cache interface {
	GetSecret(key SecretKey) (string, bool)
	PutSecret(key SecretKey, secret string)
}

type cacheImpl struct{}

func New() Cache {
	return &cacheImpl{}
}
