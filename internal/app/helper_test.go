package app

import (
	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/providers"
	providerMocks "github.com/lemoony/snipkit/mocks/provider"
)

func withProviderSnippets(snippets []model.Snippet) Option {
	return optionFunc(func(a *appImpl) {
		provider := providerMocks.Provider{}
		provider.On("GetSnippets").Return(snippets, nil)

		providerBuilder := providerMocks.ProviderBuilder{}
		providerBuilder.On("BuildProvider", mock.Anything, mock.Anything).Return([]providers.Provider{&provider}, nil)
		a.providersBuilder = &providerBuilder
	})
}

func withProviders(p ...providers.Provider) Option {
	return optionFunc(func(a *appImpl) {
		providerBuilder := providerMocks.ProviderBuilder{}
		providerBuilder.On("BuildProvider", mock.Anything, mock.Anything).Return(p, nil)
		a.providersBuilder = &providerBuilder
	})
}
