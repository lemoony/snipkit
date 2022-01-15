package pictarinesnip

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/utils/system"
	"github.com/lemoony/snipkit/internal/utils/testutil"
)

const (
	testDataDefaultLibraryPath = "testdata/userhome/Library/Containers/com.pictarine.Snip/Data/Library/Application Support/Snip/snippets"
)

func Test_GetSnippets(t *testing.T) {
	tests := []struct {
		name          string
		configFunc    func(config *Config)
		expectedCount int
	}{
		{
			name: "no tags",
			configFunc: func(config *Config) {
				config.Enabled = true
				config.LibraryPath = testDataDefaultLibraryPath
			},
			expectedCount: 2,
		},
		{
			name: "snipkit tag",
			configFunc: func(config *Config) {
				config.Enabled = true
				config.LibraryPath = testDataDefaultLibraryPath
				config.IncludeTags = []string{"snipkit"}
			},
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			system := testutil.NewTestSystem()

			config := Config{}
			tt.configFunc(&config)

			provider, err := NewProvider(WithSystem(system), WithConfig(config))
			assert.NoError(t, err)
			assert.NotNil(t, provider)

			assert.Len(t, provider.GetSnippets(), tt.expectedCount)
		})
	}
}

func Test_Info(t *testing.T) {
	system := testutil.NewTestSystem(
		system.WithUserContainersDir("testdata/userhome/Library/Containers"),
	)

	config := Config{}
	config.Enabled = true
	config.LibraryPath = testDataDefaultLibraryPath

	provider, err := NewProvider(WithSystem(system), WithConfig(config))
	assert.NoError(t, err)
	assert.NotNil(t, provider)

	info := provider.Info()
	assert.Len(t, info.Lines, 3)

	assert.Equal(t, "Pictarine Snip library path", info.Lines[0].Key)
	assert.Equal(t, testDataDefaultLibraryPath, info.Lines[0].Value)

	assert.Equal(t, "Pictarine Snip tags", info.Lines[1].Key)
	assert.Equal(t, "None", info.Lines[1].Value)

	assert.Equal(t, "Pictarine Snip total number of snippets", info.Lines[2].Key)
	assert.Equal(t, "2", info.Lines[2].Value)
}
