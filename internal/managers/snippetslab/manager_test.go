package snippetslab

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/utils/system"
	"github.com/lemoony/snipkit/internal/utils/testutil"
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			system := testutil.NewTestSystem()

			config := Config{}
			tt.configFunc(&config)

			manager, err := NewManager(WithSystem(system), WithConfig(config))
			assert.NoError(t, err)
			assert.NotNil(t, manager)

			assert.Len(t, manager.GetSnippets(), tt.expectedCount)
		})
	}
}

func Test_Info(t *testing.T) {
	system := testutil.NewTestSystem(
		system.WithUserContainersDir(testDataPreferencesWithUserDefinedLibraryPath),
	)

	config := Config{}
	config.Enabled = true
	config.LibraryPath = testDataDefaultLibraryPath

	manager, err := NewManager(WithSystem(system), WithConfig(config))
	assert.NoError(t, err)
	assert.NotNil(t, manager)

	info := manager.Info()
	assert.Len(t, info, 4)

	assert.Equal(t, "SnippetsLab preferences path", info[0].Key)
	assert.Equal(t, "no valid preferences url found", info[0].Value)

	assert.Equal(t, "SnippetsLab library path", info[1].Key)
	assert.Equal(t, testDataDefaultLibraryPath, info[1].Value)

	assert.Equal(t, "SnippetsLab tags", info[2].Key)
	assert.Equal(t, "None", info[2].Value)

	assert.Equal(t, "SnippetsLab total number of snippets", info[3].Key)
	assert.Equal(t, "2", info[3].Value)
}
