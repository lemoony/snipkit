package assistant

import (
	"emperror.dev/errors"

	"github.com/lemoony/snipkit/internal/assistant/gemini"
	"github.com/lemoony/snipkit/internal/assistant/openai"
)

type ClientProvider interface {
	GetClient(config Config) (Client, error)
}

type clientProviderImpl struct{}

func (p clientProviderImpl) GetClient(config Config) (Client, error) {
	key, err := config.ClientKey()
	if err != nil {
		return nil, err
	}

	switch key {
	case openai.Key:
		return openai.NewClient(openai.WithConfig(*config.OpenAI))
	case gemini.Key:
		return gemini.NewClient(gemini.WithConfig(*config.Gemini))
	default:
		return nil, errors.Errorf("Unsupported assistant key %s", key)
	}
}
