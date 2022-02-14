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

var suffixLanguageMap = map[string]model.Language{
	".sh":   model.LanguageBash,
	".yaml": model.LanguageYAML,
	".yml":  model.LanguageYAML,
	".md":   model.LanguageMarkdown,
	".toml": model.LanguageTOML,
}

type Manager struct {
	system      *system.System
	config      Config
	suffixRegex []*regexp.Regexp
}

// Option configures a Manager.
type Option interface {
	apply(p *Manager)
}

// optionFunc wraps a func so that it satisfies the Option interface.
type optionFunc func(manager *Manager)

func (f optionFunc) apply(manager *Manager) {
	f(manager)
}

// WithSystem sets the utils.System instance to be used by Manager.
func WithSystem(system *system.System) Option {
	return optionFunc(func(p *Manager) {
		p.system = system
	})
}

func WithConfig(config Config) Option {
	return optionFunc(func(p *Manager) {
		p.config = config
	})
}

func NewManager(options ...Option) (*Manager, error) {
	manager := &Manager{}

	for _, o := range options {
		o.apply(manager)
	}

	if !manager.config.Enabled {
		log.Debug().Msg("No fslibrary manager because it is disabled")
		return nil, nil
	}

	manager.compileSuffixRegex()

	return manager, nil
}

func (m Manager) Key() model.ManagerKey {
	return Key
}

func (m Manager) Info() []model.InfoLine {
	var lines []model.InfoLine

	lines = append(lines, model.InfoLine{
		IsError: false,
		Key:     "Filesystem library paths",
		Value:   fmt.Sprintf("[%s]", strings.Join(m.config.LibraryPath, ", ")),
	})

	lines = append(lines, model.InfoLine{
		IsError: false,
		Key:     "Filesystem library allowed suffixes",
		Value:   fmt.Sprintf("[%s]", strings.Join(m.config.SuffixRegex, ", ")),
	})

	lines = append(lines, model.InfoLine{
		IsError: false,
		Key:     "Filesystem library total number of snippets",
		Value:   fmt.Sprintf("%d", len(m.GetSnippets())),
	})

	return lines
}

func (m *Manager) GetSnippets() []model.Snippet {
	var result []model.Snippet

	for _, dir := range m.config.LibraryPath {
		entries, err := afero.ReadDir(m.system.Fs, dir)
		if err != nil {
			panic(err)
		}

		for _, entry := range entries {
			fileName := filepath.Base(entry.Name())
			filePath := filepath.Join(dir, fileName)

			if !checkSuffix(fileName, m.suffixRegex) {
				continue
			}

			snippet := model.Snippet{
				UUID:     filePath,
				TagUUIDs: []string{},
				LanguageFunc: func() model.Language {
					return languageForSuffix(filepath.Ext(fileName))
				},
				ContentFunc: func() string {
					contents := string(m.system.ReadFile(filePath))
					if m.config.HideTitleInPreview {
						contents = pruneTitleHeader(strings.NewReader(contents))
					}
					return contents
				},
			}

			if m.config.LazyOpen {
				snippet.SetTitle(fileName)
			} else {
				snippet.SetTitle(m.getSnippetName(filePath))
			}

			result = append(result, snippet)
		}
	}

	return result
}

func (m *Manager) Sync(events model.SyncEventChannel) {
	close(events)
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

func (m *Manager) compileSuffixRegex() {
	m.suffixRegex = make([]*regexp.Regexp, len(m.config.SuffixRegex))
	for i, s := range m.config.SuffixRegex {
		m.suffixRegex[i] = regexp.MustCompile(s)
	}
}

func (m *Manager) getSnippetName(filePath string) string {
	return getSnippetName(m.system, filePath)
}

func languageForSuffix(suffix string) model.Language {
	if e, ok := suffixLanguageMap[suffix]; ok {
		return e
	}
	return model.LanguageUnknown
}
