package providers

import "github.com/lemoony/snipkit/internal/model"

type Provider interface {
	Info() model.ProviderInfo
	GetSnippets() []model.Snippet
}
