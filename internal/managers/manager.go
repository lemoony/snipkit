package managers

import "github.com/lemoony/snipkit/internal/model"

type Manager interface {
	Info() model.ManagerInfo
	GetSnippets() []model.Snippet
}
