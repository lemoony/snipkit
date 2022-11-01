package config

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/managers/fslibrary"
	"github.com/lemoony/snipkit/internal/managers/githubgist"
	"github.com/lemoony/snipkit/internal/managers/pet"
	"github.com/lemoony/snipkit/internal/managers/pictarinesnip"
	"github.com/lemoony/snipkit/internal/managers/snippetslab"
)

func Test_serializeToYamlWithComment(t *testing.T) {
	var testConfig VersionWrapper
	testConfig.Version = version
	testConfig.Config.Editor = "foo-editor"
	testConfig.Config.Script.Shell = "/bin/zsh"
	testConfig.Config.Script.RemoveComments = true
	testConfig.Config.Script.ParameterMode = ParameterModeSet
	testConfig.Config.Style.Theme = "simple"
	testConfig.Config.Style.HideKeyMap = true

	testConfig.Config.FuzzySearch = true

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
		Enabled:            true,
		LibraryPath:        []string{"/path/to/file/system/library"},
		SuffixRegex:        []string{".sh"},
		LazyOpen:           true,
		HideTitleInPreview: true,
	}

	expectedConfigBytes, err := ioutil.ReadFile(testDataExampleConfig)
	assert.NoError(t, err)

	actualCfgBytes := SerializeToYamlWithComment(testConfig)
	assert.Equal(t, string(expectedConfigBytes), string(actualCfgBytes))
}
