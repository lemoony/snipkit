package fslibrary

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/phuslu/log"
	"github.com/spf13/afero"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/system"
)

const (
	maxLineNumberTitleComment = 3
)

var suffixLanguageMap = map[string]model.Language{
	".sh":   model.LanguageBash,
	".yaml": model.LanguageYAML,
	".yml":  model.LanguageYAML,
	".md":   model.LanguageMarkdown,
	".toml": model.LanguageTOML,
}

type Provider struct {
	system      *system.System
	config      Config
	suffixRegex []*regexp.Regexp
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
		log.Debug().Msg("No fslibrary provider because it is disabled")
		return nil, nil
	}

	provider.compileSuffixRegex()

	return provider, nil
}

func (p Provider) Info() model.ProviderInfo {
	var lines []model.ProviderLine

	lines = append(lines, model.ProviderLine{
		IsError: false,
		Key:     "Filesystem library paths",
		Value:   fmt.Sprintf("[%s]", strings.Join(p.config.LibraryPath, ", ")),
	})

	lines = append(lines, model.ProviderLine{
		IsError: false,
		Key:     "Filesystem library allowed suffixes",
		Value:   fmt.Sprintf("[%s]", strings.Join(p.config.SuffixRegex, ", ")),
	})

	lines = append(lines, model.ProviderLine{
		IsError: false,
		Key:     "Filesystem library total number of snippets",
		Value:   fmt.Sprintf("%d", len(p.GetSnippets())),
	})

	return model.ProviderInfo{
		Lines: lines,
	}
}

func (p *Provider) GetSnippets() []model.Snippet {
	var result []model.Snippet

	for _, dir := range p.config.LibraryPath {
		entries, err := afero.ReadDir(p.system.Fs, dir)
		if err != nil {
			panic(err)
		}

		for _, entry := range entries {
			fileName := filepath.Base(entry.Name())
			filePath := filepath.Join(dir, fileName)

			if !checkSuffix(fileName, p.suffixRegex) {
				continue
			}

			snippet := model.Snippet{
				UUID:     filePath,
				TagUUIDs: []string{},
				LanguageFunc: func() model.Language {
					return languageForSuffix(filepath.Ext(fileName))
				},
				ContentFunc: func() string {
					contents := string(p.system.ReadFile(filePath))
					if p.config.HideTitleInPreview {
						contents = pruneTitleHeader(strings.NewReader(contents))
					}
					return contents
				},
			}

			if p.config.LazyOpen {
				snippet.SetTitle(fileName)
			} else {
				snippet.SetTitle(p.getSnippetName(filePath))
			}

			result = append(result, snippet)
		}
	}

	return result
}

func checkSuffix(filename string, regexes []*regexp.Regexp) bool {
	if len(regexes) == 0 {
		return true
	}

	suffix := filepath.Ext(filename)

	for _, r := range regexes {
		if r.MatchString(suffix) {
			return true
		}
	}

	return false
}

func (p *Provider) compileSuffixRegex() {
	p.suffixRegex = make([]*regexp.Regexp, len(p.config.SuffixRegex))
	for i, s := range p.config.SuffixRegex {
		p.suffixRegex[i] = regexp.MustCompile(s)
	}
}

func (p *Provider) getSnippetName(filePath string) string {
	return getSnippetName(p.system, filePath)
}

func languageForSuffix(suffix string) model.Language {
	if e, ok := suffixLanguageMap[suffix]; ok {
		return e
	}
	return model.LanguageUnknown
}
