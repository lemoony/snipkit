package config

import (
	"path"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snippet-kit/internal/ui/uimsg"
	"github.com/lemoony/snippet-kit/internal/utils"
	mocks "github.com/lemoony/snippet-kit/mocks/ui"
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

	system := utils.NewTestSystem()

	v := viper.New()
	v.SetFs(system.Fs)
	v.SetConfigFile(cfgFilePath)

	terminal := &mocks.Terminal{}
	terminal.On("Confirm", mock.Anything).Return(true, nil)
	terminal.On("PrintMessage", mock.Anything).Return()

	s := NewService(WithSystem(system), WithViper(v), WithTerminal(terminal))

	assert.NoError(t, s.Create())
	assert.True(t, s.(serviceImpl).hasConfig())
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
	cfgFilePath := path.Join(t.TempDir(), "config.yaml")

	system := utils.NewTestSystem()

	v := viper.New()
	v.SetConfigFile(cfgFilePath)
	v.SetFs(system.Fs)

	terminal := &mocks.Terminal{}
	terminal.On("Confirm", mock.Anything).Return(true, nil)
	terminal.On("PrintMessage", mock.Anything).Return()

	s := NewService(WithSystem(system), WithViper(v), WithTerminal(terminal))

	assert.NoError(t, s.Create())
	assert.True(t, s.(serviceImpl).hasConfig())

	assert.NoError(t, s.Clean())
	assert.False(t, s.(serviceImpl).hasConfig())
}

func Test_Clean_Decline(t *testing.T) {
	cfgFilePath := path.Join(t.TempDir(), "config.yaml")

	system := utils.NewTestSystem()

	v := viper.New()
	v.SetConfigFile(cfgFilePath)
	v.SetFs(system.Fs)

	terminal := &mocks.Terminal{}
	terminal.On("Confirm", uimsg.ConfirmCreateConfigFile()).Return(true, nil)
	terminal.On("PrintMessage", mock.Anything).Return()

	s := NewService(WithSystem(system), WithViper(v), WithTerminal(terminal))

	assert.NoError(t, s.Create())
	assert.True(t, s.(serviceImpl).hasConfig())

	terminal.On("Confirm", uimsg.ConfirmDeleteConfigFile()).Return(false, nil)

	assert.NoError(t, s.Clean())
	assert.True(t, s.(serviceImpl).hasConfig())
}

func Test_Clean_NoConfig(t *testing.T) {
	cfgFilePath := path.Join(t.TempDir(), "config.yaml")

	system := utils.NewTestSystem()

	v := viper.New()
	v.SetConfigFile(cfgFilePath)
	v.SetFs(system.Fs)

	terminal := &mocks.Terminal{}
	terminal.On("PrintError", uimsg.NoConfig()).Return()

	s := NewService(WithSystem(system), WithViper(v), WithTerminal(terminal))

	assert.NoError(t, s.Clean())
}
