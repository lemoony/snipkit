package pictarinesnip

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/model"
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

			manager, err := NewManager(WithSystem(system), WithConfig(config))
			assert.NoError(t, err)
			assert.NotNil(t, manager)

			assert.Len(t, manager.GetSnippets(), tt.expectedCount)
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

	manager, err := NewManager(WithSystem(system), WithConfig(config))
	assert.NoError(t, err)
	assert.NotNil(t, manager)

	info := manager.Info()
	assert.Len(t, info, 3)

	assert.Equal(t, "Pictarine Snip library path", info[0].Key)
	assert.Equal(t, testDataDefaultLibraryPath, info[0].Value)

	assert.Equal(t, "Pictarine Snip tags", info[1].Key)
	assert.Equal(t, "None", info[1].Value)

	assert.Equal(t, "Pictarine Snip total number of snippets", info[2].Key)
	assert.Equal(t, "2", info[2].Value)
}

func Test_Sync(t *testing.T) {
	events := make(model.SyncEventChannel)
	manager := Manager{}
	manager.Sync(events)
	close(events)
}

func Test_SaveAssistantSnippet(t *testing.T) {
	assert.PanicsWithError(t, "Not implemented", func() {
		Manager{}.SaveAssistantSnippet("", "foo.sh", []byte("dummy content"))
	})
}
