package managers

import "github.com/lemoony/snipkit/internal/model"

type Manager interface {
	Key() model.ManagerKey
	Info() []model.InfoLine
	GetSnippets() []model.Snippet
	Sync(*model.SyncFeedback) bool
}
