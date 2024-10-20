package assistant

import "github.com/lemoony/snipkit/internal/model"

type ClientProvider interface {
	GetClient(key model.AssistantKey) Client
}
