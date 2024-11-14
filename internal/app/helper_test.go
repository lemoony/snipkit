package app

import (
	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snipkit/internal/managers"
	"github.com/lemoony/snipkit/internal/model"
	managerMocks "github.com/lemoony/snipkit/mocks/managers"
)

var testSnippetContent = `# ${VAR1} Name: First Output
# ${VAR1} Description: What to print on the tui first
echo "${VAR1}`

func withManagerSnippets(snippets []model.Snippet) Option {
	return optionFunc(func(a *appImpl) {
		manager := managerMocks.Manager{}
		manager.On("GetSnippets").Return(snippets, nil)

		provider := managerMocks.Provider{}
		provider.On("CreateManager", mock.Anything, mock.Anything, mock.Anything).Return([]managers.Manager{&manager}, nil)
		a.provider = &provider
	})
}

func withManager(m ...managers.Manager) Option {
	return optionFunc(func(a *appImpl) {
		providerBuilder := managerMocks.Provider{}
		providerBuilder.On("CreateManager", mock.Anything, mock.Anything, mock.Anything).Return(m, nil)
		a.provider = &providerBuilder
	})
}
