package pictarinesnip

import (
	"fmt"
	"strings"

	"github.com/phuslu/log"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/stringutil"
	"github.com/lemoony/snipkit/internal/utils/system"
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
		log.Debug().Msg("No pictarinesnip provider because it is disabled")
		return nil, nil
	}

	return provider, nil
}

func (p Provider) Info() model.ProviderInfo {
	var lines []model.ProviderLine

	lines = append(lines, model.ProviderLine{
		IsError: false,
		Key:     "Pictarine Snip library path",
		Value:   stringutil.StringOrDefault(p.config.LibraryPath, "not set"),
	})

	lines = append(lines, model.ProviderLine{
		IsError: false,
		Key:     "Pictarine Snip tags",
		Value:   stringutil.StringOrDefault(strings.Join(p.config.IncludeTags, ","), "None"),
	})

	lines = append(lines, model.ProviderLine{
		IsError: true, Key: "Pictarine Snip total number of snippets", Value: fmt.Sprintf("%d", len(p.GetSnippets())),
	})

	return model.ProviderInfo{
		Lines: lines,
	}
}

func (p *Provider) GetSnippets() []model.Snippet {
	tags := p.getValidTagUUIDs()
	return parseLibrary(p.config.LibraryPath, p.system, &tags)
}

func (p *Provider) getValidTagUUIDs() stringutil.StringSet {
	result := stringutil.StringSet{}
	for _, validTag := range p.config.IncludeTags {
		result.Add(validTag)
	}
	return result
}
