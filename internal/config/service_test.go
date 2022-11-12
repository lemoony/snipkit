package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snipkit/internal/config/testdata"
	"github.com/lemoony/snipkit/internal/managers"
	"github.com/lemoony/snipkit/internal/managers/fslibrary"
	"github.com/lemoony/snipkit/internal/managers/pet"
	"github.com/lemoony/snipkit/internal/managers/pictarinesnip"
	"github.com/lemoony/snipkit/internal/managers/snippetslab"
	"github.com/lemoony/snipkit/internal/ui/uimsg"
	"github.com/lemoony/snipkit/internal/utils/assertutil"
	"github.com/lemoony/snipkit/internal/utils/system"
	"github.com/lemoony/snipkit/internal/utils/testutil"
	"github.com/lemoony/snipkit/internal/utils/testutil/mockutil"
	mocks "github.com/lemoony/snipkit/mocks/ui"
)

func Test_LoadConfig(t *testing.T) {
	v := viper.New()
	v.SetConfigFile(testdata.ConfigPath(t, testdata.Latest))

	s := NewService(WithViper(v))

	config, err := s.LoadConfig()

	assert.NoError(t, err)
	assert.NotNil(t, config)

	assert.Equal(t, "foo-editor", config.Editor)
	assert.Equal(t, "simple", config.Style.Theme)
	assert.True(t, config.Manager.SnippetsLab.Enabled)
	assert.Equal(t, "/path/to/lib", config.Manager.SnippetsLab.LibraryPath)
	assert.Len(t, config.Manager.SnippetsLab.IncludeTags, 2)
}

func Test_Create(t *testing.T) {
	snipkitHomeDir := t.TempDir()
	_ = os.Setenv("SNIPKIT_HOME", snipkitHomeDir)
	cfgFilePath := path.Join(snipkitHomeDir, "test-config.yaml")
	system := testutil.NewTestSystem()

	v := viper.New()
	v.SetFs(system.Fs)
	v.SetConfigFile(cfgFilePath)

	confirm := uimsg.ConfigFileCreateConfirm(cfgFilePath, snipkitHomeDir, false)
	tui := &mocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	tui.On(mockutil.Confirmation, confirm, mock.Anything).Return(true, nil)
	tui.On(mockutil.Print, mock.Anything).Return()

	s := NewService(WithSystem(system), WithViper(v), WithTerminal(tui))

	s.Create()
	tui.AssertCalled(t, mockutil.Confirmation, confirm, mock.Anything)
	tui.AssertNumberOfCalls(t, mockutil.Confirmation, 1)

	assert.True(t, s.(*serviceImpl).hasConfig())
}

func Test_Create_Decline(t *testing.T) {
	snipkitHomeDir := t.TempDir()
	_ = os.Setenv("SNIPKIT_HOME", snipkitHomeDir)
	cfgFilePath := path.Join(snipkitHomeDir, "test-config.yaml")

	system := testutil.NewTestSystem()

	v := viper.New()
	v.SetFs(system.Fs)
	v.SetConfigFile(cfgFilePath)

	tui := &mocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	tui.
		On(mockutil.Confirmation, uimsg.ConfigFileCreateConfirm(cfgFilePath, snipkitHomeDir, false), mock.Anything).
		Return(false, nil)

	tui.On(mockutil.Print, mock.Anything).Return()

	s := NewService(WithSystem(system), WithViper(v), WithTerminal(tui))

	s.Create()
	assert.False(t, s.(*serviceImpl).hasConfig())
	tui.AssertCalled(
		t,
		mockutil.Confirmation,
		uimsg.ConfigFileCreateConfirm(cfgFilePath, snipkitHomeDir, false),
		mock.Anything,
	)
	tui.AssertNumberOfCalls(t, mockutil.Confirmation, 1)

	if exists, err := afero.Exists(system.Fs, cfgFilePath); err != nil {
		assert.NoError(t, err)
	} else {
		assert.False(t, exists)
	}
}

func Test_Create_Recreate_Decline(t *testing.T) {
	snipkitHomeDir := t.TempDir()
	_ = os.Setenv("SNIPKIT_HOME", snipkitHomeDir)
	cfgFilePath := path.Join(snipkitHomeDir, "test-config.yaml")

	system := testutil.NewTestSystem()
	system.WriteFile(cfgFilePath, []byte{})

	v := viper.New()
	v.SetFs(system.Fs)
	v.SetConfigFile(cfgFilePath)

	tui := &mocks.TUI{}
	tui.On(mockutil.Print, mock.Anything).Return()
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	tui.On(mockutil.Confirmation, uimsg.ConfigFileCreateConfirm(cfgFilePath, snipkitHomeDir, true), mock.Anything).
		Return(false, nil)

	s := NewService(WithSystem(system), WithViper(v), WithTerminal(tui))
	assert.True(t, s.(*serviceImpl).hasConfig())

	s.Create()
	tui.AssertCalled(
		t, mockutil.Confirmation, uimsg.ConfigFileCreateConfirm(cfgFilePath, snipkitHomeDir, true), mock.Anything,
	)
	tui.AssertNumberOfCalls(t, mockutil.Confirmation, 1)

	// assert file was not changed and, thus, is still empty
	assert.Empty(t, system.ReadFile(cfgFilePath))
}

func Test_Info(t *testing.T) {
	v := viper.New()
	v.SetConfigFile(testdata.ConfigPath(t, testdata.Latest))

	tui := &mocks.TUI{}
	tui.On("OpenEditor", testdata.ConfigPath(t, testdata.Latest), "foo-editor").Return(nil)
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()

	s := NewService(WithViper(v), WithTerminal(tui))

	infos := s.Info()

	assert.Equal(t, "Config path", infos[0].Key)
	assert.Equal(t, "SNIPKIT_HOME", infos[1].Key)
	assert.Equal(t, "Theme", infos[2].Key)
}

func Test_Edit(t *testing.T) {
	v := viper.New()
	v.SetConfigFile(testdata.ConfigPath(t, testdata.Latest))

	tui := &mocks.TUI{}
	tui.On("OpenEditor", testdata.ConfigPath(t, testdata.Latest), "foo-editor").Return(nil)
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()

	s := NewService(WithViper(v), WithTerminal(tui))

	s.Edit()
}

func Test_Clean(t *testing.T) {
	system := testutil.NewTestSystem(system.WithConfigCome("~/.snipkit"))

	// create a config file
	system.CreatePath(system.ConfigPath())
	system.WriteFile(system.ConfigPath(), []byte(""))

	// create a custom theme
	themePath := filepath.Join(system.ThemesDir(), "custom.yaml")
	system.CreatePath(themePath)
	system.WriteFile(themePath, []byte(""))

	v := viper.New()
	v.SetConfigFile(system.ConfigPath())
	v.SetFs(system.Fs)

	tui := &mocks.TUI{}
	tui.On(mockutil.Confirmation, mock.Anything, mock.Anything).Return(true, nil)
	tui.On(mockutil.Print, mock.Anything).Return()
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()

	s := NewService(WithSystem(system), WithViper(v), WithTerminal(tui))

	assert.True(t, s.(*serviceImpl).hasConfig())

	s.Clean()

	tui.AssertCalled(
		t, mockutil.Confirmation, uimsg.ConfigFileDeleteConfirm(v.ConfigFileUsed()), mock.Anything,
	)
	tui.AssertCalled(t, mockutil.Confirmation, uimsg.ThemesDeleteConfirm(system.ThemesDir()), mock.Anything)

	assert.False(t, s.(*serviceImpl).hasConfig())
	assertutil.AssertExists(t, system.Fs, system.ConfigPath(), false)
	assertutil.AssertExists(t, system.Fs, system.ThemesDir(), false)
}

func Test_Clean_Decline(t *testing.T) {
	system := testutil.NewTestSystem(system.WithConfigCome("~/.snipkit"))

	// create a config file
	system.CreatePath(system.ConfigPath())
	system.WriteFile(system.ConfigPath(), []byte(""))

	// create a custom theme
	themePath := filepath.Join(system.ThemesDir(), "custom.yaml")
	system.CreatePath(themePath)
	system.WriteFile(themePath, []byte(""))

	v := viper.New()
	v.SetConfigFile(system.ConfigPath())
	v.SetFs(system.Fs)

	tui := &mocks.TUI{}

	tui.On(mockutil.Print, mock.Anything).Return()
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	tui.
		On(mockutil.Confirmation, uimsg.ConfigFileDeleteConfirm(system.ConfigPath()), mock.Anything).
		Return(false, nil)
	tui.
		On(mockutil.Confirmation, uimsg.ThemesDeleteConfirm(system.ThemesDir()), mock.Anything).
		Return(false, nil)

	s := NewService(WithSystem(system), WithViper(v), WithTerminal(tui))
	assert.True(t, s.(*serviceImpl).hasConfig())

	s.Clean()

	tui.AssertCalled(
		t, mockutil.Confirmation, uimsg.ConfigFileDeleteConfirm(system.ConfigPath()), mock.Anything,
	)
	tui.AssertCalled(t, mockutil.Confirmation, uimsg.ThemesDeleteConfirm(system.ThemesDir()), mock.Anything)

	assert.True(t, s.(*serviceImpl).hasConfig())
	assertutil.AssertExists(t, system.Fs, system.ConfigPath(), true)
}

func Test_Clean_NoConfig(t *testing.T) {
	cfgFilePath := path.Join(t.TempDir(), "config.yaml")

	system := testutil.NewTestSystem(system.WithConfigCome(cfgFilePath))

	tui := mocks.TUI{}
	tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
	tui.On(mockutil.Print, mock.Anything).Return()

	v := viper.New()
	v.SetConfigFile(cfgFilePath)
	v.SetFs(system.Fs)

	s := NewService(WithSystem(system), WithTerminal(&tui), WithViper(v))
	s.Clean()

	tui.AssertCalled(t, mockutil.Print, uimsg.ConfigNotFound(cfgFilePath))
}

func Test_ConfigFilePath(t *testing.T) {
	cfgFilePath := path.Join(t.TempDir(), "config.yaml")

	system := testutil.NewTestSystem()

	v := viper.New()
	v.SetConfigFile(cfgFilePath)
	v.SetFs(system.Fs)

	s := NewService(WithSystem(system), WithViper(v))

	assert.Equal(t, cfgFilePath, s.ConfigFilePath())
}

func Test_UpdateManagerConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		update managers.Config
		assert func(cfg Config)
	}{
		{
			name: "snippetslab", update: managers.Config{SnippetsLab: &snippetslab.Config{Enabled: true}},
			assert: func(cfg Config) { assert.True(t, cfg.Manager.SnippetsLab.Enabled) },
		},
		{
			name: "pictarinesnip", update: managers.Config{PictarineSnip: &pictarinesnip.Config{Enabled: true}},
			assert: func(cfg Config) { assert.True(t, cfg.Manager.PictarineSnip.Enabled) },
		},
		{
			name: "üet", update: managers.Config{Pet: &pet.Config{Enabled: true}},
			assert: func(cfg Config) { assert.True(t, cfg.Manager.Pet.Enabled) },
		},
		{
			name: "fslibrary", update: managers.Config{FsLibrary: &fslibrary.Config{Enabled: true}}, assert: func(cfg Config) {
				assert.True(t, cfg.Manager.FsLibrary.Enabled)
			},
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfgPath := filepath.Join(t.TempDir(), "cfg.yaml")
			s := testutil.NewTestSystem()

			v := viper.New()
			v.SetConfigFile(cfgPath)
			v.SetFs(s.Fs)

			createConfigFile(s, v)
			service := NewService(WithSystem(s), WithViper(v))
			if cfg, err := service.LoadConfig(); err != nil {
				assert.NoError(t, err)
			} else {
				assert.Nil(t, cfg.Manager.SnippetsLab)
				assert.Nil(t, cfg.Manager.PictarineSnip)
				assert.Nil(t, cfg.Manager.FsLibrary)
			}

			service.UpdateManagerConfig(tt.update)

			// create new service instance so that old config isn't cached anymore
			service = NewService(WithSystem(s), WithViper(v))
			cfg, err := service.LoadConfig()
			assert.NoError(t, err)
			tt.assert(cfg)
		})
	}
}

func Test_NeedsMigration(t *testing.T) {
	v := viper.New()
	v.SetConfigFile(testdata.ConfigPath(t, testdata.ConfigV110))

	s := NewService(WithViper(v))

	needsMigration, version := s.NeedsMigration()

	assert.True(t, needsMigration)
	assert.Equal(t, "1.1.0", version)
}

func Test_Migrate(t *testing.T) {
	tests := []struct {
		confirm bool
	}{
		{confirm: true},
		{confirm: false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Migrate - confirm: %t", tt.confirm), func(t *testing.T) {
			config := testdata.ConfigBytes(t, testdata.Latest)
			config110 := testdata.ConfigBytes(t, testdata.ConfigV110)

			s := testutil.NewTestSystem()

			cfgFilePath := filepath.Join(s.UserHome(), "config.yml")
			s.WriteFile(cfgFilePath, config110)

			v := viper.New()
			v.SetFs(s.Fs)
			v.SetConfigFile(cfgFilePath)

			tui := mocks.TUI{}
			tui.On(mockutil.ApplyConfig, mock.Anything, mock.Anything).Return()
			tui.On(mockutil.Print, mock.Anything).Return()
			tui.
				On(mockutil.Confirmation, uimsg.ConfigFileMigrationConfirm(string(config)), mock.Anything).
				Return(tt.confirm)

			service := NewService(WithViper(v), WithTerminal(&tui), WithSystem(s))
			service.Migrate()

			tui.AssertCalled(t, mockutil.Print, uimsg.ConfigFileMigrationResult(tt.confirm, cfgFilePath))

			// check if config is migrated or not
			if tt.confirm {
				assert.Equal(t, config, s.ReadFile(cfgFilePath))
			} else {
				assert.Equal(t, config110, s.ReadFile(cfgFilePath))
			}
		})
	}
}
