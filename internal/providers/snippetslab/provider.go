package snippetslab

import (
	"fmt"
	"strings"

	"github.com/phuslu/log"

	"github.com/lemoony/snippet-kit/internal/model"
	"github.com/lemoony/snippet-kit/internal/utils/stringutil"
	"github.com/lemoony/snippet-kit/internal/utils/system"
)

type Provider struct {
	system *system.System
	config Config
}

// Option configures a Provider.
type Option interface {
	apply(p *Provider)
}

// optionFunc wraps a func so that it satisfies the Option interface.
type optionFunc func(provider *Provider)

func (f optionFunc) apply(provider *Provider) {
	f(provider)
}

// WithSystem sets the utils.System instance to be used by Provider.
func WithSystem(system *system.System) Option {
	return optionFunc(func(p *Provider) {
		p.system = system
	})
}

func WithConfig(config Config) Option {
	return optionFunc(func(p *Provider) {
		p.config = config
	})
}

func NewProvider(options ...Option) (*Provider, error) {
	provider := &Provider{}

	for _, o := range options {
		o.apply(provider)
	}

	if !provider.config.Enabled {
		log.Debug().Msg("No snippetsLab provider because it is disabled")
		return nil, nil
	}

	return provider, nil
}

func (p Provider) libraryPath() snippetsLabLibrary {
	return snippetsLabLibrary(p.config.LibraryPath)
}

func (p Provider) Info() model.ProviderInfo {
	var lines []model.ProviderLine

	if path, err := findPreferencesPath(p.system); err != nil {
		lines = append(lines, model.ProviderLine{IsError: true, Key: "SnippetsLab preferences path", Value: err.Error()})
	} else {
		lines = append(lines, model.ProviderLine{
			IsError: true, Key: "SnippetsLab preferences path", Value: path,
		})
	}

	lines = append(lines, model.ProviderLine{
		IsError: false, Key: "SnippetsLab library path", Value: p.config.LibraryPath,
	})

	lines = append(lines, model.ProviderLine{
		IsError: false,
		Key:     "SnippetsLab tags",
		Value:   stringutil.StringOrDefault(strings.Join(p.config.IncludeTags, ","), "None"),
	})

	lines = append(lines, model.ProviderLine{
		IsError: true, Key: "SnippetsLab total number of snippets", Value: fmt.Sprintf("%d", len(p.GetSnippets())),
	})

	return model.ProviderInfo{
		Lines: lines,
	}
}

func (p *Provider) GetSnippets() []model.Snippet {
	validTagUUIDs := p.getValidTagUUIDs()

	snippets, err := parseSnippets(p.libraryPath())
	if err != nil {
		panic(err)
	}

	if len(validTagUUIDs) == 0 {
		return snippets
	} else {
		var result []model.Snippet
		for _, snippet := range snippets {
			if hasValidTag(snippet.TagUUIDs, validTagUUIDs) {
				result = append(result, snippet)
			}
		}
		return result
	}
}

func hasValidTag(snippetTagUUIDS []string, validTagUUIDs stringutil.StringSet) bool {
	for _, tagUUID := range snippetTagUUIDS {
		if validTagUUIDs.Contains(tagUUID) {
			return true
		}
	}
	return false
}

func (p *Provider) getValidTagUUIDs() stringutil.StringSet {
	tags, err := parseTags(p.libraryPath())
	if err != nil {
		panic(err)
	}

	result := stringutil.StringSet{}
	for _, validTag := range p.config.IncludeTags {
		for tagKey, tagValue := range tags {
			if strings.Compare(tagValue, validTag) == 0 {
				result.Add(tagKey)
			}
		}
	}

	return result
}
