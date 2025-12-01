package assistant

import (
	"time"

	"emperror.dev/errors"
	"github.com/phuslu/log"

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
	descriptions := make([]model.AssistantDescription, 0, len(SupportedProviders))

	for _, providerType := range SupportedProviders {
		info := GetProviderInfo(providerType)
		enabled := false

		// Check if this provider is currently enabled in config
		for _, p := range config.Providers {
			if p.Type == providerType && p.Enabled {
				enabled = true
				break
			}
		}

		descriptions = append(descriptions, model.AssistantDescription{
			Key:         model.AssistantKey(providerType),
			Name:        info.Name,
			Description: info.Description,
			Enabled:     enabled,
		})
	}

	return descriptions
}

func (a *assistantImpl) AutoConfig(key model.AssistantKey) Config {
	result := a.config
	providerType := ProviderType(key)

	// Disable all existing providers
	for i := range result.Providers {
		result.Providers[i].Enabled = false
	}

	// Find existing provider or add new one
	found := false
	for i := range result.Providers {
		if result.Providers[i].Type == providerType {
			result.Providers[i].Enabled = true
			found = true
			break
		}
	}

	if !found {
		// Add new provider with defaults
		newProvider := DefaultProviderConfig(providerType)
		result.Providers = append(result.Providers, newProvider)
	}

	return result
}
