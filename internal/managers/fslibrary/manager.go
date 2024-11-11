package fslibrary

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/afero"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui"
	"github.com/lemoony/snipkit/internal/ui/uimsg"
	"github.com/lemoony/snipkit/internal/utils/idutil"
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
	printer     ui.MessagePrinter
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

func WithPrinter(printer ui.MessagePrinter) Option {
	return optionFunc(func(p *Manager) {
		p.printer = printer
	})
}

func NewManager(options ...Option) (*Manager, error) {
	manager := &Manager{}
	for _, o := range options {
		o.apply(manager)
	}
	manager.compileSuffixRegex()
	return manager, nil
}

func (m Manager) Key() model.ManagerKey {
	return Key
}

func (m Manager) SaveAssistantSnippet(snippetTitle string, filename string, contents []byte) {
	dirPath := m.config.LibraryPath[m.config.AssistantLibraryPathIndex]
	if file, err := filepath.Abs(filepath.Join(dirPath, filename)); err == nil {
		m.system.CreatePath(file)
		m.system.WriteFile(file, []byte(formatSnippet(string(contents), snippetTitle)))
		m.printer.Print(uimsg.AssistantSnippetSaved(snippetTitle, file))
	} else {
		panic(err)
	}
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
		result = append(result, m.snippetsFromDir(dir)...)
	}
	return result
}

func (m *Manager) Sync(model.SyncEventChannel) {
	// do nothing
}

func (m *Manager) snippetsFromDir(dir string) []model.Snippet {
	var result []model.Snippet

	entries, err := afero.ReadDir(m.system.Fs, dir)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			result = append(result, m.snippetsFromDir(path.Join(dir, entry.Name()))...)
			continue
		}

		fileName := filepath.Base(entry.Name())
		filePath := filepath.Join(dir, fileName)

		if !checkSuffix(fileName, m.suffixRegex) {
			continue
		}

		snippet := snippetImpl{
			id:   idutil.FormatSnippetID(filePath, idPrefix),
			path: filePath,
			tags: []string{},
			contentFunc: func() string {
				contents := string(m.system.ReadFile(filePath))
				if m.config.HideTitleInPreview {
					contents = pruneTitleHeader(strings.NewReader(contents))
				}
				return contents
			},
			titleFunc: func() string {
				if m.config.LazyOpen {
					return fileName
				} else {
					return m.getSnippetName(filePath)
				}
			},
		}

		result = append(result, &snippet)
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

func (m *Manager) compileSuffixRegex() {
	m.suffixRegex = make([]*regexp.Regexp, len(m.config.SuffixRegex))
	for i, s := range m.config.SuffixRegex {
		m.suffixRegex[i] = regexp.MustCompile(s)
	}
}

func (m *Manager) getSnippetName(filePath string) string {
	return getSnippetName(m.system, filePath)
}

func LanguageForSuffix(suffix string) model.Language {
	if e, ok := suffixLanguageMap[suffix]; ok {
		return e
	}
	return model.LanguageUnknown
}
