package managers

import "github.com/lemoony/snipkit/internal/model"

type Manager interface {
	Key() model.ManagerKey
	Info() []model.InfoLine
	GetSnippets() []model.Snippet
	Sync(model.SyncEventChannel)
	SaveAssistantSnippet(snippetTitle string, filename string, contents []byte)
}
