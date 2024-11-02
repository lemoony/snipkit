package masscode

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/system"
	"github.com/lemoony/snipkit/internal/utils/testutil"
)

func Test_GetInfo(t *testing.T) {
	config := Config{
		Enabled:      true,
		MassCodeHome: testDataMassCodeV2Path,
		Version:      version2,
	}

	sys := testutil.NewTestSystem(system.WithUserHome(testDataUserHomeV2))
	provider, err := NewManager(WithSystem(sys), WithConfig(config))
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
		name        string
		userHome    string
		version     Version
		tags        []string
		expectedLen int
	}{
		{name: "v2 - no tags", userHome: testDataUserHomeV2, version: version2, tags: []string{}, expectedLen: 3},
		{name: "v2 - 1 tag", userHome: testDataUserHomeV2, version: version2, tags: []string{"snipkit"}, expectedLen: 1},
		{name: "v2 - tag excludes all", userHome: testDataUserHomeV2, version: version2, tags: []string{"foo"}, expectedLen: 0},
		{name: "v1 - no tags", userHome: testDataUserHomeV1, version: version1, tags: []string{}, expectedLen: 2},
		{name: "v1 - 1 tags", userHome: testDataUserHomeV1, version: version1, tags: []string{"snipkit"}, expectedLen: 1},
		{name: "v1 - tage excludes all", userHome: testDataUserHomeV1, version: version1, tags: []string{"foo"}, expectedLen: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sys := testutil.NewTestSystem(system.WithUserHome(tt.userHome))
			config := Config{
				Enabled:      true,
				MassCodeHome: filepath.Join(tt.userHome, defaultMassCodeHomePath),
				Version:      tt.version,
				IncludeTags:  tt.tags,
			}

			provider, _ := NewManager(WithSystem(sys), WithConfig(config))
			assert.Len(t, provider.GetSnippets(), tt.expectedLen)
		})
	}
}

func Test_SaveAssistantSnippet(t *testing.T) {
	assert.PanicsWithError(t, "Not implemented", func() {
		Manager{}.SaveAssistantSnippet("", "foo.sh", []byte("dummy content"))
	})
}
