package pet

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/testutil"
)

func Test_GetInfo(t *testing.T) {
	config := Config{
		Enabled:      true,
		LibraryPaths: []string{testDataSnippetFile},
	}

	system := testutil.NewTestSystem()
	provider, err := NewManager(WithSystem(system), WithConfig(config))
	assert.NoError(t, err)

	info := provider.Info()

	assert.Len(t, info, 3)

	assert.Equal(t, info[0].Key, "Pet enabled")
	assert.Equal(t, info[0].Value, "true")
	assert.False(t, info[0].IsError)

	assert.Equal(t, info[1].Key, "Pet snippet file paths")
	assert.Equal(t, info[1].Value, testDataSnippetFile)
	assert.False(t, info[1].IsError)

	assert.Equal(t, info[2].Key, "Pet total number of snippets")
	assert.Equal(t, info[2].Value, "2")
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
		includeTags              []string
		expectedNumberOfSnippets int
	}{
		{name: "no include tags", includeTags: []string{}, expectedNumberOfSnippets: 2},
		{name: "include single tag", includeTags: []string{"tag1"}, expectedNumberOfSnippets: 1},
		{name: "exclude all ", includeTags: []string{"unknown"}, expectedNumberOfSnippets: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{
				Enabled:      true,
				LibraryPaths: []string{testDataSnippetFile},
				IncludeTags:  tt.includeTags,
			}

			system := testutil.NewTestSystem()
			provider, _ := NewManager(WithSystem(system), WithConfig(config))
			assert.Len(t, provider.GetSnippets(), tt.expectedNumberOfSnippets)
		})
	}
}

func Test_SaveAssistantSnippet(t *testing.T) {
	assert.PanicsWithError(t, "Not implemented", func() {
		Manager{}.SaveAssistantSnippet("", "foo.sh", []byte("dummy content"))
	})
}
