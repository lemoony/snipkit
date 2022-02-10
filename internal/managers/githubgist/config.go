package githubgist

import (
	"fmt"
	"regexp"

	"emperror.dev/errors"

	"github.com/lemoony/snipkit/internal/utils/system"
)

type (
	AuthMethod      string
	SnippetNameMode string
)

const (
	AuthMethodNone  = AuthMethod("None")
	AuthMethodToken = AuthMethod("Token")
	AuthMethodOAuth = AuthMethod("OAuth")

	apiURLPattern = "https://api.%s/users/%s/gists"

	SnippetNameModeDescription = "DESCRIPTION"
	SnippetNameModeFilename    = "FILENAME"
)

var urlRegex = regexp.MustCompile("^gist.(.*)/(.*)$")

type Config struct {
	Enabled bool         `yaml:"enabled" head_comment:"If set to false, github gist is disabled completely."`
	Gists   []GistConfig `yaml:"gists" head_comment:"You can define multiple independent Github Gist sources."`
}

type GistConfig struct {
	Enabled                   bool            `yaml:"enabled" head_comment:"If set to false, this github gist url is ignored."`
	URL                       string          `yaml:"url" head_comment:"URL to the GitHub gist account."`
	AuthenticationMethod      AuthMethod      `yaml:"authenticationMethod" head_comment:"Supported values: None, OAuth, Token. Default value: None (which means no authentication). In order to retrieve secret gists, you must be authenticated."`
	IncludeTags               []string        `yaml:"includeTags" head_comment:"If this list is not empty, only those gists that match the listed tags will be provided to you."`
	SuffixRegex               []string        `yaml:"suffixRegex" head_comment:"Only gist files with endings which match one of the listed suffixes will be considered."`
	NameMode                  SnippetNameMode `yaml:"nameMode" head_comment:"Defines where the snippet name is extracted from (see also titleHeaderEnabled). Allowed values: DESCRIPTION, FILENAME."`
	RemoveTagsFromDescription bool            `yaml:"removeTagsFromDescription" head_comment:"If set to true, any tags will be removed from the description."`
	TitleHeaderEnabled        bool            `yaml:"titleHeaderEnabled" head_comment:"If set to true, the snippet title can be overwritten by defining a title header within the gist."`
}

func (g GistConfig) apiURL() string {
	matches := urlRegex.FindStringSubmatch(g.URL)
	const minMatches = 3
	if len(matches) < minMatches {
		panic(errors.Errorf("invalid gist url: %s", g.URL))
	}
	return fmt.Sprintf(apiURLPattern, matches[1], matches[2])
}

func AutoDiscoveryConfig(system *system.System) *Config {
	return &Config{
		Enabled: false,
		Gists: []GistConfig{
			{
				Enabled:              false,
				URL:                  "gist.github.com/<USERNAME>",
				AuthenticationMethod: AuthMethodNone,
				IncludeTags:          []string{},
			},
		},
	}
}
