package assistant

import (
	"github.com/lemoony/snipkit/internal/assistant/openai"
	"github.com/lemoony/snipkit/internal/cache"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/system"
)

type Provider interface {
	Query(string) string
	AssistantDescriptions(Config) []model.AssistantDescription
	AutoConfig(model.AssistantKey, *system.System) Config
}

type providerImpl struct {
	system *system.System
	config Config
	cache  cache.Cache
}

func NewBuilder(system *system.System, config Config, cache cache.Cache) Provider {
	return providerImpl{system: system, config: config, cache: cache}
}

func (a providerImpl) Query(prompt string) string {
	client, err := openai.NewClient(openai.WithCache(a.cache))
	if err != nil {
		panic(err)
	}

	response := client.Query(prompt)
	return extractBashScript(response)
}

func (a providerImpl) AutoConfig(key model.AssistantKey, s *system.System) Config {
	return Config{}
}

func (a providerImpl) AssistantDescriptions(config Config) []model.AssistantDescription {
	return nil
}
