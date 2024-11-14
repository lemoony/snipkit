package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/assistant"
	"github.com/lemoony/snipkit/internal/assistant/openai"
	"github.com/lemoony/snipkit/internal/config/testdata"
	"github.com/lemoony/snipkit/internal/managers/fslibrary"
	"github.com/lemoony/snipkit/internal/managers/githubgist"
	"github.com/lemoony/snipkit/internal/managers/pet"
	"github.com/lemoony/snipkit/internal/managers/pictarinesnip"
	"github.com/lemoony/snipkit/internal/managers/snippetslab"
)

//nolint:funlen // test function for yaml config is allowed to be too long
func Test_serializeToYamlWithComment(t *testing.T) {
	var testConfig VersionWrapper
	testConfig.Version = Version
	testConfig.Config.Editor = "foo-editor"
	testConfig.Config.Script.Shell = "/bin/zsh"
	testConfig.Config.Script.RemoveComments = true
	testConfig.Config.Script.ParameterMode = ParameterModeSet
	testConfig.Config.Script.ExecConfirm = false
	testConfig.Config.Script.ExecPrint = false
	testConfig.Config.Style.Theme = "simple"
	testConfig.Config.Style.HideKeyMap = true

	testConfig.Config.FuzzySearch = true
	testConfig.Config.SecretStorage = SecretStorageKeyring

	testConfig.Config.Assistant = assistant.Config{
		SaveMode: assistant.SaveModeNever,
		OpenAI: &openai.Config{
			Enabled:   true,
			Endpoint:  "test.endpoint.com",
			Model:     "test/model",
			APIKeyEnv: "foo-key",
		},
	}

	testConfig.Config.Manager.SnippetsLab = &snippetslab.Config{
		Enabled: true, LibraryPath: "/path/to/lib", IncludeTags: []string{"snipkit", "othertag"},
	}
	testConfig.Config.Manager.PictarineSnip = &pictarinesnip.Config{
		Enabled: false, LibraryPath: "", IncludeTags: []string{},
	}
	testConfig.Config.Manager.Pet = &pet.Config{
		Enabled:      true,
		LibraryPaths: []string{"/foouser/.config/pet/snippet.toml"},
	}
	testConfig.Config.Manager.GithubGist = &githubgist.Config{
		Enabled: true,
		Gists: []githubgist.GistConfig{
			{
				Enabled:                   true,
				URL:                       "gist.github.com/<yourUser>",
				AuthenticationMethod:      githubgist.AuthMethodPAT,
				IncludeTags:               []string{},
				SuffixRegex:               []string{},
				NameMode:                  githubgist.SnippetNameModeCombinePreferDescription,
				TitleHeaderEnabled:        true,
				RemoveTagsFromDescription: true,
				HideTitleInPreview:        true,
			},
		},
	}
	testConfig.Config.Manager.FsLibrary = &fslibrary.Config{
		Enabled:                   true,
		LibraryPath:               []string{"/path/to/file/system/library"},
		AssistantLibraryPathIndex: 0,
		SuffixRegex:               []string{".sh"},
		LazyOpen:                  true,
		HideTitleInPreview:        true,
	}

	expectedConfigBytes := testdata.ConfigBytes(t, testdata.Example)
	actualCfgBytes := SerializeToYamlWithComment(testConfig)
	assert.Equal(t, string(expectedConfigBytes), string(actualCfgBytes))
}
