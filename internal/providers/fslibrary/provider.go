package fslibrary

import (
	"bufio"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/phuslu/log"
	"github.com/spf13/afero"

	"github.com/lemoony/snippet-kit/internal/model"
	"github.com/lemoony/snippet-kit/internal/utils/system"
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

	if s, err := p.GetSnippets(); err != nil {
		lines = append(lines, model.ProviderLine{
			IsError: true,
			Key:     "Filesystem library total number of snippets",
			Value:   err.Error(),
		})
	} else {
		lines = append(lines, model.ProviderLine{
			IsError: false,
			Key:     "Filesystem library total number of snippets",
			Value:   fmt.Sprintf("%d", len(s)),
		})
	}

	return model.ProviderInfo{
		Lines: lines,
	}
}

func (p *Provider) GetSnippets() ([]model.Snippet, error) {
	var result []model.Snippet

	for _, dir := range p.config.LibraryPath {
		entries, err := afero.ReadDir(p.system.Fs, dir)
		if err != nil {
			return []model.Snippet{}, err
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
					return string(p.system.ReadFile(filePath))
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

	return result, nil
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
	file, err := p.system.Fs.Open(filePath)
	fileName := filepath.Base(filePath)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	titleLine := 0

	title := ""
	for scanner.Scan() {
		lineNumber++

		if lineNumber-titleLine > maxLineNumberTitleComment {
			return fileName
		}

		line := scanner.Text()

		if strings.HasPrefix(line, "#") {
			switch {
			case titleLine == 0 && len(strings.TrimSpace(line)) == 1:
				titleLine++
			case titleLine == 1:
				title = strings.TrimSpace(strings.TrimPrefix(line, "#"))
				titleLine++
			case titleLine == 2 && len(strings.TrimSpace(line)) == 1:
				if title != "" {
					return title
				} else {
					return fileName
				}
			}
		}
	}

	return fileName
}

func languageForSuffix(suffix string) model.Language {
	if e, ok := suffixLanguageMap[suffix]; ok {
		return e
	}
	return model.LanguageUnknown
}
