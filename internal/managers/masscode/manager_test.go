package masscode

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/model"
	system "github.com/lemoony/snipkit/internal/utils/system"
	"github.com/lemoony/snipkit/internal/utils/testutil"
)

func Test_GetInfo(t *testing.T) {
	config := Config{
		Enabled:      true,
		MassCodeHome: testDataMassCodeV2Path,
		Version:      version2,
	}

	system := testutil.NewTestSystem()
	provider, err := NewManager(WithSystem(system), WithConfig(config))
	assert.NoError(t, err)

	info := provider.Info()

	assert.Len(t, info, 3)

	assert.Equal(t, info[0].Key, "MassCode enabled")
	assert.Equal(t, "true", info[0].Value)
	assert.False(t, info[0].IsError)

	assert.Equal(t, info[1].Key, "MassCode path")
	assert.Equal(t, testDataMassCodeV2Path, info[1].Value)
	assert.False(t, info[1].IsError)

	assert.Equal(t, info[2].Key, "MassCode total number of Snippets")
	assert.Equal(t, "3", info[2].Value)
	assert.False(t, info[2].IsError)
}

func Test_Key(t *testing.T) {
	assert.Equal(t, Key, Manager{}.Key())
}

func Test_Sync(t *testing.T) {
	events := make(model.SyncEventChannel)
	manager := Manager{}
	manager.Sync(events)
	close(events)
}

func Test_GetSnippets(t *testing.T) {
	tests := []struct {
		name                     string
		userHome                 string
		version                  Version
		includeTags              []string
		expectedNumberOfSnippets int
	}{
		{
			name:                     "v2 - no tags",
			userHome:                 testDataUserHomeV2,
			version:                  version2,
			includeTags:              []string{},
			expectedNumberOfSnippets: 3,
		},
		{
			name:                     "v2 - 1 tag",
			userHome:                 testDataUserHomeV2,
			version:                  version2,
			includeTags:              []string{"snipkit"},
			expectedNumberOfSnippets: 1,
		},
		{
			name:                     "v2 - tag excludes all",
			userHome:                 testDataUserHomeV2,
			version:                  version2,
			includeTags:              []string{"foo"},
			expectedNumberOfSnippets: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sys := testutil.NewTestSystem(system.WithUserHome(tt.userHome))
			config := Config{
				Enabled:      true,
				MassCodeHome: filepath.Join(tt.userHome, defaultMassCodeHomePath),
				Version:      tt.version,
				IncludeTags:  tt.includeTags,
			}

			provider, _ := NewManager(WithSystem(sys), WithConfig(config))
			assert.Len(t, provider.GetSnippets(), tt.expectedNumberOfSnippets)
		})
	}
}
