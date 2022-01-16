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
	"github.com/lemoony/snipkit/internal/utils/testutil/mockutil"
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
	assert.True(t, config.Manager.SnippetsLab.Enabled)
	assert.Equal(t, "/path/to/lib", config.Manager.SnippetsLab.LibraryPath)
	assert.Len(t, config.Manager.SnippetsLab.IncludeTags, 2)
}

func Test_Create(t *testing.T) {
	cfgFilePath := path.Join(t.TempDir(), "test-config.yaml")

	system := testutil.NewTestSystem()

	v := viper.New()
	v.SetFs(system.Fs)
	v.SetConfigFile(cfgFilePath)

	confirm := uimsg.ConfigFileCreateConfirm(cfgFilePath, "", false)
	terminal := &mocks.Terminal{}
	terminal.On(mockutil.Confirmation, confirm, mock.Anything).Return(true, nil)
	terminal.On(mockutil.PrintMessage, mock.Anything).Return()

	s := NewService(WithSystem(system), WithViper(v), WithTerminal(terminal))

	s.Create()
	terminal.AssertCalled(t, mockutil.Confirmation, confirm, mock.Anything)
	terminal.AssertNumberOfCalls(t, mockutil.Confirmation, 1)

	assert.True(t, s.(serviceImpl).hasConfig())
}

func Test_Create_Decline(t *testing.T) {
	cfgFilePath := path.Join(t.TempDir(), "test-config.yaml")

	system := testutil.NewTestSystem()

	v := viper.New()
	v.SetFs(system.Fs)
	v.SetConfigFile(cfgFilePath)

	terminal := &mocks.Terminal{}
	terminal.
		On(mockutil.Confirmation, uimsg.ConfigFileCreateConfirm(cfgFilePath, "", false), mock.Anything).
		Return(false, nil)

	terminal.On(mockutil.PrintMessage, mock.Anything).Return()

	s := NewService(WithSystem(system), WithViper(v), WithTerminal(terminal))

	s.Create()
	assert.False(t, s.(serviceImpl).hasConfig())
	terminal.AssertCalled(
		t,
		mockutil.Confirmation,
		uimsg.ConfigFileCreateConfirm(cfgFilePath, "", false),
		mock.Anything,
	)
	terminal.AssertNumberOfCalls(t, mockutil.Confirmation, 1)

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
	terminal.
		On(mockutil.Confirmation, uimsg.ConfigFileCreateConfirm(cfgFilePath, "", false), mock.Anything).
		Return(true, nil)
	terminal.
		On(mockutil.Confirmation, uimsg.ConfigFileCreateConfirm(cfgFilePath, "", true), mock.Anything).
		Return(false, nil)
	terminal.
		On(mockutil.PrintMessage, mock.Anything).
		Return()

	s := NewService(WithSystem(system), WithViper(v), WithTerminal(terminal))

	s.Create()
	assert.True(t, s.(serviceImpl).hasConfig())
	terminal.AssertCalled(
		t, mockutil.Confirmation, uimsg.ConfigFileCreateConfirm(cfgFilePath, "", false), mock.Anything,
	)
	terminal.AssertNumberOfCalls(t, mockutil.Confirmation, 1)

	stat, _ := system.Fs.Stat(cfgFilePath)
	modTime := stat.ModTime()

	s.Create()

	assert.True(t, s.(serviceImpl).hasConfig())
	terminal.AssertCalled(
		t, mockutil.Confirmation, uimsg.ConfigFileCreateConfirm(cfgFilePath, "", true), mock.Anything,
	)
	terminal.AssertNumberOfCalls(t, mockutil.Confirmation, 2)

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
	terminal.On(mockutil.Confirmation, mock.Anything, mock.Anything).Return(true, nil)
	terminal.On(mockutil.PrintMessage, mock.Anything).Return()

	s := NewService(WithSystem(system), WithViper(v), WithTerminal(terminal))

	assert.True(t, s.(serviceImpl).hasConfig())

	s.Clean()

	terminal.AssertCalled(
		t, mockutil.Confirmation, uimsg.ConfigFileDeleteConfirm(v.ConfigFileUsed()), mock.Anything,
	)
	terminal.AssertCalled(t, mockutil.Confirmation, uimsg.ThemesDeleteConfirm(system.ThemesDir()), mock.Anything)

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

	terminal.On(mockutil.PrintMessage, mock.Anything).Return()
	terminal.
		On(mockutil.Confirmation, uimsg.ConfigFileDeleteConfirm(system.ConfigPath()), mock.Anything).
		Return(false, nil)
	terminal.
		On(mockutil.Confirmation, uimsg.ThemesDeleteConfirm(system.ThemesDir()), mock.Anything).
		Return(false, nil)

	s.Clean()

	terminal.AssertCalled(
		t, mockutil.Confirmation, uimsg.ConfigFileDeleteConfirm(system.ConfigPath()), mock.Anything,
	)
	terminal.AssertCalled(t, mockutil.Confirmation, uimsg.ThemesDeleteConfirm(system.ThemesDir()), mock.Anything)

	assert.True(t, s.(serviceImpl).hasConfig())
	assertutil.AssertExists(t, system.Fs, system.ConfigPath(), true)
}

func Test_Clean_NoConfig(t *testing.T) {
	cfgFilePath := path.Join(t.TempDir(), "config.yaml")

	system := testutil.NewTestSystem(system.WithConfigCome(cfgFilePath))

	term := mocks.Terminal{}
	term.On(mockutil.PrintMessage, mock.Anything).Return()

	v := viper.New()
	v.SetConfigFile(cfgFilePath)
	v.SetFs(system.Fs)

	s := NewService(WithSystem(system), WithTerminal(&term), WithViper(v))
	s.Clean()

	term.AssertCalled(t, mockutil.PrintMessage, uimsg.ConfigNotFound(cfgFilePath))
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

func Test_initConfigHelpText(t *testing.T) {
}
