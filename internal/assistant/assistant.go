package assistant

import (
	"github.com/lemoony/snipkit/internal/assistant/gemini"
	"github.com/lemoony/snipkit/internal/assistant/openai"
	"github.com/lemoony/snipkit/internal/cache"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/system"
)

type Assistant interface {
	Query(string) (string, string)
	AutoConfig(model.AssistantKey) Config
	AssistantDescriptions(config Config) []model.AssistantDescription
}

type assistantImpl struct {
	system   *system.System
	config   Config
	cache    cache.Cache
	provider ClientProvider
}

func NewBuilder(system *system.System, config Config, cache cache.Cache, options ...Option) Assistant {
	asst := assistantImpl{system: system, config: config, cache: cache, provider: clientProviderImpl{}}
	for _, o := range options {
		o.apply(&asst)
	}
	return asst
}

func (a assistantImpl) Query(prompt string) (string, string) {
	client, err := a.provider.GetClient(a.config)
	if err != nil {
		panic(err)
	}

	response, err := client.Query(prompt)
	if err != nil {
		panic(err)
	}
	return extractBashScript(response)
}

func (a assistantImpl) AssistantDescriptions(config Config) []model.AssistantDescription {
	return []model.AssistantDescription{
		openai.Description(config.OpenAI),
		gemini.Description(config.Gemini),
	}
}

func (a assistantImpl) AutoConfig(key model.AssistantKey) Config {
	result := a.config
	if key == openai.Key {
		updated := openai.AutoDiscoveryConfig(a.config.OpenAI)
		result.OpenAI = &updated
		if result.Gemini != nil {
			result.Gemini.Enabled = false
		}
	} else if key == gemini.Key {
		updated := gemini.AutoDiscoveryConfig(a.config.Gemini)
		result.Gemini = &updated
		if result.OpenAI != nil {
			result.OpenAI.Enabled = false
		}
	}
	return result
}
