package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_serializeToYamlWithComment(t *testing.T) {
	var testConfig versionWrapper
	testConfig.Version = "1.0.0"
	testConfig.Config.Editor = "code"
	testConfig.Config.Style.Theme = "dracula"
	testConfig.Config.Providers.SnippetsLab.Enabled = true
	testConfig.Config.Providers.SnippetsLab.LibraryPath = "/path/to/lib"
	testConfig.Config.Providers.SnippetsLab.IncludeTags = []string{"snipkit", "othertag"}
	testConfig.Config.Providers.SnippetsLab.ExcludeTags = []string{}

	exepectedYaml := `version: 1.0.0
config:
    style:
        # The theme defines the terminal colors used by Snipkit.
        # Available themes:dracula,solarized,light,dark.
        theme: dracula
    # Your preferred editor to open the config file when typing 'snipkit config edit'.
    editor: code # Defaults to a reasonable value for your operation system when empty.
    provider:
        snippetsLab:
            # Set to true if you want to use SnippetsLab.
            enabled: true
            # Path to your *.snippetslablibrary file.
            # SnipKit will try to detect this file automatically when generating the config.
            libraryPath: /path/to/lib
            # If this list is not empty, only those snippets that match the listed tags will be provided to you.
            includeTags:
                - snipkit
                - othertag
            # If this list is not empty, snippets that have one of the listed tags will not be provided to you.
            excludeTags: []
        fileSystem:
            # If set to false, the files specified via libraryPath will not be provided to you.
            enabled: false
            # Path to a directory that holds snippets files. Each file must hold one snippet only.
            libraryPath: ""
`

	yaml, err := serializeToYamlWithComment(testConfig)
	assert.NoError(t, err)
	assert.Equal(t, exepectedYaml, string(yaml))
}
