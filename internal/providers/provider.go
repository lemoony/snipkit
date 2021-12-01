package providers

import "github.com/lemoony/snippet-kit/internal/model"

type Provider interface {
	Info() model.ProviderInfo
	GetSnippets() ([]model.Snippet, error)
}
