package assistant

import (
	"emperror.dev/errors"

	assistErrors "github.com/lemoony/snipkit/internal/assistant/errors"
	"github.com/lemoony/snipkit/internal/assistant/gemini"
	"github.com/lemoony/snipkit/internal/assistant/openai"
	"github.com/lemoony/snipkit/internal/cache"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/system"
)

type Assistant interface {
	Query(string) string
	AutoConfig(model.AssistantKey, *system.System) Config
	AssistantDescriptions(config Config) []model.AssistantDescription
}

type assistantImpl struct {
	system *system.System
	config Config
	cache  cache.Cache
}

func NewBuilder(system *system.System, config Config, cache cache.Cache) Assistant {
	return assistantImpl{system: system, config: config, cache: cache}
}

func (a assistantImpl) Query(prompt string) string {
	client, err := a.getClient()
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

func (a assistantImpl) AutoConfig(key model.AssistantKey, s *system.System) Config {
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

func (a assistantImpl) getClient() (Client, error) {
	switch {
	case a.config.moreThanOneEnabled():
		panic(errors.New("Invalid config: more than one assistant is enabled."))
	case a.config.OpenAI != nil && a.config.OpenAI.Enabled:
		return openai.NewClient(openai.WithConfig(*a.config.OpenAI))
	case a.config.Gemini != nil && a.config.Gemini.Enabled:
		return gemini.NewClient(gemini.WithConfig(*a.config.Gemini))
	}
	return nil, assistErrors.ErrorNoClientConfiguredOrEnabled
}
