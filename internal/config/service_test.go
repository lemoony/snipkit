package config

import (
	"path"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snipkit/internal/ui/uimsg"
	"github.com/lemoony/snipkit/internal/utils/assertutil"
	"github.com/lemoony/snipkit/internal/utils/system"
	"github.com/lemoony/snipkit/internal/utils/testutil"
	mocks "github.com/lemoony/snipkit/mocks/ui"
)

const testDataExampleConfig = "testdata/example-config.yaml"

func Test_LoadConfig(t *testing.T) {
	v := viper.New()
	v.SetConfigFile(testDataExampleConfig)

	s := NewService(WithViper(v))

	config, err := s.LoadConfig()

	assert.NoError(t, err)
	assert.NotNil(t, config)

	assert.Equal(t, "foo-editor", config.Editor)
	assert.Equal(t, "dracula", config.Style.Theme)
	assert.True(t, config.Providers.SnippetsLab.Enabled)
	assert.Equal(t, "/path/to/lib", config.Providers.SnippetsLab.LibraryPath)
	assert.Len(t, config.Providers.SnippetsLab.IncludeTags, 2)
	assert.Len(t, config.Providers.SnippetsLab.ExcludeTags, 0)
}

func Test_Create(t *testing.T) {
	cfgFilePath := path.Join(t.TempDir(), "test-config.yaml")

	system := testutil.NewTestSystem()

	v := viper.New()
	v.SetFs(system.Fs)
	v.SetConfigFile(cfgFilePath)

	terminal := &mocks.Terminal{}
	terminal.On("Confirm", uimsg.ConfirmCreateConfigFile(cfgFilePath)).Return(true, nil)
	terminal.On("PrintMessage", mock.Anything).Return()

	s := NewService(WithSystem(system), WithViper(v), WithTerminal(terminal))

	s.Create()
	terminal.AssertCalled(t, "Confirm", uimsg.ConfirmCreateConfigFile(cfgFilePath))
	terminal.AssertNumberOfCalls(t, "Confirm", 1)

	assert.True(t, s.(serviceImpl).hasConfig())
}

func Test_Create_Decline(t *testing.T) {
	cfgFilePath := path.Join(t.TempDir(), "test-config.yaml")

	system := testutil.NewTestSystem()

	v := viper.New()
	v.SetFs(system.Fs)
	v.SetConfigFile(cfgFilePath)

	terminal := &mocks.Terminal{}
	terminal.On("Confirm", uimsg.ConfirmCreateConfigFile(cfgFilePath)).Return(false, nil)

	terminal.On("PrintMessage", mock.Anything).Return()

	s := NewService(WithSystem(system), WithViper(v), WithTerminal(terminal))

	s.Create()
	assert.False(t, s.(serviceImpl).hasConfig())
	terminal.AssertCalled(t, "Confirm", uimsg.ConfirmCreateConfigFile(cfgFilePath))
	terminal.AssertNumberOfCalls(t, "Confirm", 1)

	if exists, err := afero.Exists(system.Fs, cfgFilePath); err != nil {
		assert.NoError(t, err)
	} else {
		assert.False(t, exists)
	}
}

func Test_Create_Recreate_Decline(t *testing.T) {
	cfgFilePath := path.Join(t.TempDir(), "test-config.yaml")

	system := testutil.NewTestSystem()

	v := viper.New()
	v.SetFs(system.Fs)
	v.SetConfigFile(cfgFilePath)

	terminal := &mocks.Terminal{}
	terminal.On("Confirm", uimsg.ConfirmCreateConfigFile(cfgFilePath)).Return(true, nil)
	terminal.On("Confirm", uimsg.ConfirmRecreateConfigFile(cfgFilePath)).Return(false, nil)

	terminal.On("PrintMessage", mock.Anything).Return()

	s := NewService(WithSystem(system), WithViper(v), WithTerminal(terminal))

	s.Create()
	assert.True(t, s.(serviceImpl).hasConfig())
	terminal.AssertCalled(t, "Confirm", uimsg.ConfirmCreateConfigFile(cfgFilePath))
	terminal.AssertNumberOfCalls(t, "Confirm", 1)

	stat, _ := system.Fs.Stat(cfgFilePath)
	modTime := stat.ModTime()

	s.Create()

	assert.True(t, s.(serviceImpl).hasConfig())
	terminal.AssertCalled(t, "Confirm", uimsg.ConfirmRecreateConfigFile(cfgFilePath))
	terminal.AssertNumberOfCalls(t, "Confirm", 2)

	// assert file was not changed by comparing the last modification time
	stat, _ = system.Fs.Stat(cfgFilePath)
	assert.Equal(t, modTime, stat.ModTime())
}

func Test_Edit(t *testing.T) {
	v := viper.New()
	v.SetConfigFile(testDataExampleConfig)

	terminal := &mocks.Terminal{}
	terminal.On("OpenEditor", testDataExampleConfig, "foo-editor").Return(nil)

	s := NewService(WithViper(v), WithTerminal(terminal))

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

	terminal := &mocks.Terminal{}
	terminal.On("Confirm", mock.Anything).Return(true, nil)
	terminal.On("PrintMessage", mock.Anything).Return()

	s := NewService(WithSystem(system), WithViper(v), WithTerminal(terminal))

	assert.True(t, s.(serviceImpl).hasConfig())

	s.Clean()

	terminal.AssertCalled(t, "Confirm", uimsg.ConfirmDeleteConfigFile(system.ConfigPath()))
	terminal.AssertCalled(t, "PrintMessage", uimsg.ConfigFileDeleted(system.ConfigPath()))

	terminal.AssertCalled(t, "Confirm", uimsg.ConfirmDeleteThemesDir(system.ThemesDir()))
	terminal.AssertCalled(t, "PrintMessage", uimsg.ThemesDeleted())

	assert.False(t, s.(serviceImpl).hasConfig())
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

	terminal := &mocks.Terminal{}

	s := NewService(WithSystem(system), WithViper(v), WithTerminal(terminal))
	assert.True(t, s.(serviceImpl).hasConfig())

	terminal.On("PrintMessage", mock.Anything).Return()
	terminal.On("Confirm", uimsg.ConfirmDeleteConfigFile(system.ConfigPath())).Return(false, nil)
	terminal.On("Confirm", uimsg.ConfirmDeleteThemesDir(system.ThemesDir())).Return(false, nil)

	s.Clean()

	terminal.AssertCalled(t, "Confirm", uimsg.ConfirmDeleteConfigFile(system.ConfigPath()))
	terminal.AssertCalled(t, "PrintMessage", uimsg.ConfigNotDeleted())

	terminal.AssertCalled(t, "Confirm", uimsg.ConfirmDeleteThemesDir(system.ThemesDir()))
	terminal.AssertCalled(t, "PrintMessage", uimsg.ThemesNotDeleted())

	assert.True(t, s.(serviceImpl).hasConfig())
	assertutil.AssertExists(t, system.Fs, system.ConfigPath(), true)
}

func Test_Clean_NoConfig(t *testing.T) {
	cfgFilePath := path.Join(t.TempDir(), "config.yaml")

	system := testutil.NewTestSystem(system.WithConfigCome(cfgFilePath))

	term := mocks.Terminal{}
	term.On("PrintMessage", mock.Anything).Return()

	v := viper.New()
	v.SetConfigFile(cfgFilePath)
	v.SetFs(system.Fs)

	s := NewService(WithSystem(system), WithTerminal(&term), WithViper(v))
	s.Clean()

	term.AssertCalled(t, "PrintMessage", uimsg.ConfigNotFound(cfgFilePath))
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
