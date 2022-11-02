package config

import (
	"io/ioutil"
	"path/filepath"

	"emperror.dev/errors"
	"github.com/phuslu/log"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/lemoony/snipkit/internal/config/migrations"
	"github.com/lemoony/snipkit/internal/managers"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/ui"
	"github.com/lemoony/snipkit/internal/ui/uimsg"
	"github.com/lemoony/snipkit/internal/utils/stringutil"
	"github.com/lemoony/snipkit/internal/utils/system"
)

var invalidConfig = Config{}

// NewService creates a new Service.
func NewService(options ...Option) Service {
	service := serviceImpl{
		v:      viper.GetViper(),
		system: system.NewSystem(),
		tui:    ui.NewTUI(),
	}
	for _, o := range options {
		o.apply(&service)
	}

	service.applyConfigTheme()

	return &service
}

type Service interface {
	Create()
	LoadConfig() (Config, error)
	Edit()
	Clean()
	UpdateManagerConfig(config managers.Config)
	NeedsMigration() (bool, string, string)
	Migrate(bool) string
	ConfigFilePath() string
	Info() []model.InfoLine
}

type serviceImpl struct {
	v       *viper.Viper
	system  *system.System
	tui     ui.TUI
	version string
	config  *Config
}

func (s *serviceImpl) Create() {
	recreate := s.hasConfig()
	confirmed := s.tui.Confirmation(
		uimsg.ConfigFileCreateConfirm(s.v.ConfigFileUsed(), s.system.HomeEnvValue(), recreate),
	)

	if confirmed {
		createConfigFile(s.system, s.v)
	}

	s.tui.Print(uimsg.ConfigFileCreateResult(confirmed, s.v.ConfigFileUsed(), recreate))
}

func (s *serviceImpl) LoadConfig() (Config, error) {
	if s.config != nil {
		return *s.config, nil
	}

	log.Debug().Msgf("SnipKit Home: %s", s.system.HomeDir())

	if !s.hasConfig() {
		return invalidConfig, ErrConfigNotFound{s.v.ConfigFileUsed()}
	}

	// If a config file is found, read it in.
	if err := s.v.ReadInConfig(); err == nil {
		log.Debug().Str("config file", s.v.ConfigFileUsed())
	} else {
		return invalidConfig, errors.Wrap(ErrInvalidConfig, "failed to read config")
	}

	var wrapper VersionWrapper
	if err := s.v.Unmarshal(&wrapper); err != nil {
		return invalidConfig, err
	}

	if wrapper.Version != Version {
		log.Warn().Msgf("Config Version is not up to date - expected %s, actual %s", Version, wrapper.Version)
		s.version = wrapper.Version
	}

	s.config = &wrapper.Config

	return *s.config, nil
}

func (s *serviceImpl) Edit() {
	cfgEditor := ""
	if cfg, err := s.LoadConfig(); errors.Is(err, ErrConfigNotFound{}) {
		panic(err)
	} else {
		cfgEditor = cfg.Editor
	}

	s.tui.OpenEditor(s.v.ConfigFileUsed(), cfgEditor)
}

func (s *serviceImpl) Clean() {
	configPath := s.v.ConfigFileUsed()
	s.applyConfigTheme()

	if s.hasConfig() {
		confirmed := s.tui.Confirmation(uimsg.ConfigFileDeleteConfirm(configPath))
		if confirmed {
			s.system.Remove(s.v.ConfigFileUsed())
		}
		s.tui.Print(uimsg.ConfigFileDeleteResult(confirmed, s.v.ConfigFileUsed()))
	} else {
		s.tui.Print(uimsg.ConfigNotFound(configPath))
	}

	if s.hasThemes() {
		confirmed := s.tui.Confirmation(uimsg.ThemesDeleteConfirm(s.system.ThemesDir()))
		if confirmed {
			s.system.RemoveAll(s.system.ThemesDir())
		}
		s.tui.Print(uimsg.ThemesDeleteResult(confirmed, s.system.ThemesDir()))
	}

	s.deleteDirectoryIfEmpty(s.system.ThemesDir())
	s.deleteDirectoryIfEmpty(filepath.Dir(s.system.ConfigPath()))

	if exists, _ := afero.DirExists(s.system.Fs, s.system.HomeDir()); exists {
		s.tui.Print(uimsg.HomeDirectoryStillExists(s.system.HomeDir()))
	}
}

func (s *serviceImpl) ConfigFilePath() string {
	return s.v.ConfigFileUsed()
}

func (s *serviceImpl) UpdateManagerConfig(managerConfig managers.Config) {
	config, err := s.LoadConfig()
	if err != nil {
		panic(errors.Wrapf(ErrInvalidConfig, "failed to load config: %s", err.Error()))
	}

	if cfg := managerConfig.FsLibrary; cfg != nil {
		config.Manager.FsLibrary = cfg
	}
	if cfg := managerConfig.SnippetsLab; cfg != nil {
		config.Manager.SnippetsLab = cfg
	}
	if cfg := managerConfig.PictarineSnip; cfg != nil {
		config.Manager.PictarineSnip = cfg
	}
	if cfg := managerConfig.Pet; cfg != nil {
		config.Manager.Pet = cfg
	}
	if cfg := managerConfig.MassCode; cfg != nil {
		config.Manager.MassCode = cfg
	}
	if cfg := managerConfig.GithubGist; cfg != nil {
		config.Manager.GithubGist = cfg
	}

	bytes := SerializeToYamlWithComment(wrap(config))
	s.system.WriteFile(s.ConfigFilePath(), bytes)
}

func (s *serviceImpl) NeedsMigration() (bool, string, string) {
	return s.version != Version, s.version, Version
}

func (s *serviceImpl) Migrate(migrate bool) string {
	if s.version == Version {
		panic(errors.Errorf("Config is already up to date: %s", Version))
	}

	result := s.updateConfigToLatest()

	return string(SerializeToYamlWithComment(result))
}

func (s *serviceImpl) hasConfig() bool {
	ok, _ := afero.Exists(s.system.Fs, s.v.ConfigFileUsed())
	return ok
}

func (s *serviceImpl) hasThemes() bool {
	themesDir := s.system.ThemesDir()
	if exists, _ := afero.DirExists(s.system.Fs, themesDir); !exists {
		return false
	}
	return !s.system.IsEmpty(themesDir)
}

func (s *serviceImpl) deleteDirectoryIfEmpty(path string) {
	if s.system.DirExists(path) && s.system.IsEmpty(path) {
		s.system.Remove(path)
	}
}

func (s *serviceImpl) applyConfigTheme() {
	cfg, err := s.LoadConfig()
	if err == nil {
		ui.ApplyConfig(cfg.Style, s.system)
	} else {
		ui.ApplyConfig(ui.DefaultConfig(), s.system)
	}
}

func (s *serviceImpl) Info() []model.InfoLine {
	result := []model.InfoLine{
		{Key: "Config path", Value: s.ConfigFilePath()},
		{Key: "SNIPKIT_HOME", Value: stringutil.StringOrDefault(s.system.HomeEnvValue(), "Not set")},
	}

	cfg, err := s.LoadConfig()
	if err == nil {
		result = append(result, model.InfoLine{Key: "Theme", Value: cfg.Style.Theme})
	}

	return result
}

func (s *serviceImpl) updateConfigToLatest() VersionWrapper {
	configBytes, err := ioutil.ReadFile(s.v.ConfigFileUsed())
	if err != nil {
		panic(err)
	}
	newConfig := migrations.Migrate(configBytes)

	var result VersionWrapper
	if err := yaml.Unmarshal(newConfig, &result); err != nil {
		panic(err)
	}

	return result
}
