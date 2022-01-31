package managers

import "github.com/lemoony/snipkit/internal/model"

type Manager interface {
	Info() []model.InfoLine
	GetSnippets() []model.Snippet
	Sync() model.SyncResult
}
