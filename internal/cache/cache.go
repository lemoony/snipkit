package cache

type Service string

const (
	defaultLabel = "snipkit"

	ServicePrivateAccessToken = Service("GitHub Private Access Token")
	ServiceOAuthAccessToken   = Service("GitHub OAuth Access Token")
)

type Cache interface {
	GetSecret(service Service) ([]byte, bool)
	PutSecret(service Service, secret []byte)
}

type cacheImpl struct {
	// label is only used for testing purposes to overwrite the defaultLabel
	label string
}

func New() Cache {
	return &cacheImpl{label: defaultLabel}
}
