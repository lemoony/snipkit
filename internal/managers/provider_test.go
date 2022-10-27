package managers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/cache"
	"github.com/lemoony/snipkit/internal/managers/fslibrary"
	"github.com/lemoony/snipkit/internal/managers/githubgist"
	"github.com/lemoony/snipkit/internal/managers/masscode"
	"github.com/lemoony/snipkit/internal/managers/pet"
	"github.com/lemoony/snipkit/internal/managers/pictarinesnip"
	"github.com/lemoony/snipkit/internal/managers/snippetslab"
	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/testutil"
)

func Test_CreateManager(t *testing.T) {
	tests := testConfigs()
	system := testutil.NewTestSystem()
	for _, tt := range tests {
		t.Run(string(tt.key), func(t *testing.T) {
			config := Config{}
			tt.configFunc(&config)

			provider := NewBuilder(cache.New(system))
			managers := provider.CreateManager(*system, config)
			assert.Len(t, managers, 1)
		})
	}
}

func Test_ManagerDescriptions(t *testing.T) {
	tests := testConfigs()
	system := testutil.NewTestSystem()
	for _, tt := range tests {
		t.Run(string(tt.key), func(t *testing.T) {
			config := Config{}
			tt.configFunc(&config)

			provider := NewBuilder(cache.New(system))
			managers := provider.ManagerDescriptions(config)
			found := false
			for _, manager := range managers {
				if manager.Key == tt.key {
					found = true
				}
			}

			assert.False(t, found)
		})
	}
}

func Test_AutoConfig(t *testing.T) {
	tests := testConfigs()
	system := testutil.NewTestSystem()

	for _, tt := range tests {
		t.Run(string(tt.key), func(t *testing.T) {
			provider := NewBuilder(cache.New(system))
			config := provider.AutoConfig(tt.key, system)

			switch tt.key {
			case snippetslab.Key:
				assert.NotNil(t, config.SnippetsLab)
			case pictarinesnip.Key:
				assert.NotNil(t, config.PictarineSnip)
			case pet.Key:
				assert.NotNil(t, config.Pet)
			case masscode.Key:
				assert.NotNil(t, config.MassCode)
			case githubgist.Key:
				assert.NotNil(t, config.GithubGist)
			case fslibrary.Key:
				assert.NotNil(t, config.FsLibrary)
			}
		})
	}
}

func testConfigs() []struct {
	key        model.ManagerKey
	configFunc func(config *Config)
} {
	return []struct {
		key        model.ManagerKey
		configFunc func(config *Config)
	}{
		{
			key: snippetslab.Key,
			configFunc: func(config *Config) {
				config.SnippetsLab = &snippetslab.Config{
					Enabled: true,
				}
			},
		},
		{
			key: pictarinesnip.Key,
			configFunc: func(config *Config) {
				config.PictarineSnip = &pictarinesnip.Config{
					Enabled: true,
				}
			},
		},
		{
			key: pet.Key,
			configFunc: func(config *Config) {
				config.Pet = &pet.Config{
					Enabled: true,
				}
			},
		},
		{
			key: masscode.Key,
			configFunc: func(config *Config) {
				config.MassCode = &masscode.Config{
					Enabled: true,
				}
			},
		},
		{
			key: githubgist.Key,
			configFunc: func(config *Config) {
				config.GithubGist = &githubgist.Config{
					Enabled: true,
				}
			},
		},
		{
			key: fslibrary.Key,
			configFunc: func(config *Config) {
				config.FsLibrary = &fslibrary.Config{
					Enabled: true,
				}
			},
		},
	}
}
