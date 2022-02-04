package config

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/managers/fslibrary"
	"github.com/lemoony/snipkit/internal/managers/githubgist"
	"github.com/lemoony/snipkit/internal/managers/pictarinesnip"
	"github.com/lemoony/snipkit/internal/managers/snippetslab"
)

func Test_serializeToYamlWithComment(t *testing.T) {
	var testConfig VersionWrapper
	testConfig.Version = "1.0.0"
	testConfig.Config.Editor = "foo-editor"
	testConfig.Config.Style.Theme = "simple"
	testConfig.Config.Manager.SnippetsLab = &snippetslab.Config{}
	testConfig.Config.Manager.SnippetsLab.Enabled = true
	testConfig.Config.Manager.SnippetsLab.LibraryPath = "/path/to/lib"
	testConfig.Config.Manager.SnippetsLab.IncludeTags = []string{"snipkit", "othertag"}

	testConfig.Config.Manager.PictarineSnip = &pictarinesnip.Config{}
	testConfig.Config.Manager.PictarineSnip.Enabled = false
	testConfig.Config.Manager.PictarineSnip.LibraryPath = ""
	testConfig.Config.Manager.PictarineSnip.IncludeTags = []string{}

	testConfig.Config.Manager.GithubGist = &githubgist.Config{}
	testConfig.Config.Manager.GithubGist.Enabled = true
	testConfig.Config.Manager.GithubGist.Gists = []githubgist.GistConfig{
		{
			Enabled:              true,
			Host:                 "github.com",
			Username:             "<yourUser>",
			AuthenticationMethod: githubgist.AuthMethodToken,
			IncludeTags:          []string{},
			SuffixRegex:          []string{},
		},
	}

	testConfig.Config.Manager.FsLibrary = &fslibrary.Config{}
	testConfig.Config.Manager.FsLibrary.Enabled = true
	testConfig.Config.Manager.FsLibrary.LibraryPath = []string{"/path/to/file/system/library"}
	testConfig.Config.Manager.FsLibrary.SuffixRegex = []string{".sh"}
	testConfig.Config.Manager.FsLibrary.LazyOpen = true
	testConfig.Config.Manager.FsLibrary.HideTitleInPreview = true

	expectedConfigBytes, err := ioutil.ReadFile(testDataExampleConfig)
	assert.NoError(t, err)

	actualCfgBytes := SerializeToYamlWithComment(testConfig)
	assert.Equal(t, string(expectedConfigBytes), string(actualCfgBytes))
}
