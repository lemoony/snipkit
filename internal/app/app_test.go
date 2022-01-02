package app

import (
	"errors"
	"io/ioutil"
	"path"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/yaml.v3"

	"github.com/lemoony/snippet-kit/internal/config"
	"github.com/lemoony/snippet-kit/internal/model"
	"github.com/lemoony/snippet-kit/internal/providers"
	"github.com/lemoony/snippet-kit/mocks"
)

func Test_NewApp_NoConfigFile(t *testing.T) {
	v := viper.NewWithOptions()

	_, err := NewApp(v)
	assert.True(t, errors.Is(err, config.ErrNoConfigFound))
}

func Test_NewAppInvalidConfigFile(t *testing.T) {
	cfgFile := path.Join(t.TempDir(), "invalid-config")
	assert.NoError(t, ioutil.WriteFile(cfgFile, []byte("invalid"), 0o600))

	v := viper.NewWithOptions()
	v.SetConfigFile(cfgFile)

	_, err := NewApp(v)
	assert.True(t, errors.Is(err, config.ErrInvalidConfig))
}

func Test_NewAppNoProviders(t *testing.T) {
	cfg := config.VersionWrapper{Version: "1.0.0", Config: config.Config{}}
	cfgFile := path.Join(t.TempDir(), "temp-config.yaml")

	if cfgBytes, err := yaml.Marshal(cfg); err != nil {
		assert.NoError(t, err)
	} else {
		assert.NoError(t, ioutil.WriteFile(cfgFile, cfgBytes, 0o600))
	}

	v := viper.NewWithOptions()
	v.SetConfigFile(cfgFile)

	term := mocks.Terminal{}
	term.On("ApplyConfig", mock.AnythingOfType("ui.Config")).Return()

	builder := mocks.Builder{}
	builder.On("BuildProvider", mock.Anything, mock.Anything).Return([]providers.Provider{}, nil)

	app, err := NewApp(v, WithTerminal(&term), WithProvidersBuilder(&builder))

	assert.NoError(t, err)

	term.AssertNumberOfCalls(t, "ApplyConfig", 1)
	builder.AssertNumberOfCalls(t, "BuildProvider", 1)

	assert.Len(t, app.Providers, 0)
}

func Test_App_GetAllSnippets(t *testing.T) {
	snippets := []model.Snippet{
		{UUID: "uuid1", Title: "title-1", Language: model.LanguageYAML, TagUUIDs: []string{}, Content: "content-1"},
		{UUID: "uuid2", Title: "title-2", Language: model.LanguageBash, TagUUIDs: []string{}, Content: "content-2"},
	}

	provider := mocks.Provider{}

	provider.On("GetSnippets").Return(snippets, nil)
	app := App{Providers: []providers.Provider{&provider}}

	s, err := app.GetAllSnippets()
	assert.NoError(t, err)
	assert.Equal(t, snippets, s)
}
