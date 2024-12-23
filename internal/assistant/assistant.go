package assistant

import (
	"time"

	"emperror.dev/errors"
	"github.com/phuslu/log"

	"github.com/lemoony/snipkit/internal/assistant/gemini"
	"github.com/lemoony/snipkit/internal/assistant/openai"
	"github.com/lemoony/snipkit/internal/cache"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui/uimsg"
	"github.com/lemoony/snipkit/internal/utils/system"
)

type Assistant interface {
	Initialize() (bool, uimsg.Printable)
	Query(string) ParsedScript
	AutoConfig(model.AssistantKey) Config
	AssistantDescriptions(config Config) []model.AssistantDescription
}

type assistantImpl struct {
	system *system.System
	config Config
	cache  cache.Cache

	client   Client
	provider ClientProvider

	demo            DemoConfig
	demoScriptIndex int
}

func NewBuilder(system *system.System, config Config, cache cache.Cache, options ...Option) Assistant {
	asst := assistantImpl{system: system, config: config, cache: cache, provider: clientProviderImpl{}}
	for _, o := range options {
		o.apply(&asst)
	}
	return &asst
}

func (a *assistantImpl) Initialize() (bool, uimsg.Printable) {
	if c, err := a.provider.GetClient(a.config); errors.Is(err, ErrorNoAssistantEnabled) {
		return false, uimsg.AssistantNoneEnabled()
	} else if err != nil {
		panic(err)
	} else {
		a.client = c
	}
	return true, uimsg.Printable{}
}

func (a *assistantImpl) Query(prompt string) ParsedScript {
	var response string
	if len(a.demo.ScriptPaths) > 0 {
		demoScript := a.system.ReadFile(a.demo.ScriptPaths[a.demoScriptIndex])
		a.demoScriptIndex++
		time.Sleep(a.demo.QueryDuration)
		response = string(demoScript)
	} else {
		var err error
		if response, err = a.client.Query(prompt); err != nil {
			panic(err)
		}
	}

	result := parseScript(response)
	log.Trace().
		Str("filename", result.Filename).
		Str("title", result.Title).
		Str("contents", result.Contents).
		Msg("Assistant generated script")
	return result
}

func (a *assistantImpl) AssistantDescriptions(config Config) []model.AssistantDescription {
	return []model.AssistantDescription{
		openai.Description(config.OpenAI),
		gemini.Description(config.Gemini),
	}
}

func (a *assistantImpl) AutoConfig(key model.AssistantKey) Config {
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
