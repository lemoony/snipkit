package assistant

import (
	"github.com/lemoony/snipkit/internal/assistant/langchain"
)

type ClientProvider interface {
	GetClient(config Config) (Client, error)
}

type clientProviderImpl struct{}

func (p clientProviderImpl) GetClient(config Config) (Client, error) {
	provider, err := config.GetActiveProvider()
	if err != nil {
		return nil, err
	}

	// Convert assistant.ProviderConfig to langchain.Config
	lcConfig := langchain.Config{
		Type:      langchain.ProviderType(provider.Type),
		Model:     provider.Model,
		APIKeyEnv: provider.APIKeyEnv,
		Endpoint:  provider.Endpoint,
		ServerURL: provider.ServerURL,
	}

	return langchain.NewClient(lcConfig)
}
