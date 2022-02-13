package githubgist

import (
	"fmt"
	"regexp"

	"emperror.dev/errors"
)

type (
	AuthMethod      string
	SnippetNameMode string
)

const (
	AuthMethodNone            = AuthMethod("None")
	AuthMethodToken           = AuthMethod("Token")
	AuthMethodOAuthDeviceFlow = AuthMethod("OAuthDeviceFlow")

	apiURLPattern  = "https://api.%s/users/%s/gists"
	hostURLPattern = "https://%s"

	SnippetNameModeDescription              = "DESCRIPTION"
	SnippetNameModeFilename                 = "FILENAME"
	SnippetNameModeCombine                  = "COMBINE"
	SnippetNameModeCombinePreferDescription = "COMBINE_PREFER_DESCRIPTION"
	SnippetNameModeCombinePreferFilename    = "COMBINE_PREFER_FILENAME"
)

var urlRegex = regexp.MustCompile("^gist.(.*)/(.*)$")

type Config struct {
	Enabled bool         `yaml:"enabled" head_comment:"If set to false, github gist is disabled completely."`
	Gists   []GistConfig `yaml:"gists" head_comment:"You can define multiple independent Github Gist sources."`
}

type GistConfig struct {
	Enabled                   bool            `yaml:"enabled" head_comment:"If set to false, this github gist url is ignored."`
	URL                       string          `yaml:"url" head_comment:"URL to the GitHub gist account."`
	AuthenticationMethod      AuthMethod      `yaml:"authenticationMethod" head_comment:"Supported values: None, OAuthDeviceFlow, Token. Default value: None (which means no authentication). In order to retrieve secret gists, you must be authenticated."`
	OAuthClientID             string          `yaml:"oauthClientID,omitempty" head_comment:"OAuth application ID, only required if the host is not github.com AND you're using OAuthDeviceFlow'."`
	IncludeTags               []string        `yaml:"includeTags" head_comment:"If this list is not empty, only those gists that match the listed tags will be provided to you."`
	SuffixRegex               []string        `yaml:"suffixRegex" head_comment:"Only gist files with endings which match one of the listed suffixes will be considered."`
	NameMode                  SnippetNameMode `yaml:"nameMode" head_comment:"Defines where the snippet name is extracted from (see also titleHeaderEnabled). Allowed values: DESCRIPTION, FILENAME, COMBINE, COMBINE_PREFER_DESCRIPTION, COMBINE_PREFER_FILENAME."`
	RemoveTagsFromDescription bool            `yaml:"removeTagsFromDescription" head_comment:"If set to true, any tags will be removed from the description."`
	TitleHeaderEnabled        bool            `yaml:"titleHeaderEnabled" head_comment:"If set to true, the snippet title can be overwritten by defining a title header within the gist."`
	HideTitleInPreview        bool            `yaml:"hideTitleInPreview" head_comment:"If set to true, the title header comment will not be shown in the preview window."`
}

func (g GistConfig) apiURL() string {
	matches := urlRegex.FindStringSubmatch(g.URL)
	const minMatches = 3
	if len(matches) < minMatches {
		panic(errors.Errorf("invalid gist url: %s", g.URL))
	}
	return fmt.Sprintf(apiURLPattern, matches[1], matches[2])
}

func (g GistConfig) hostURL() string {
	matches := urlRegex.FindStringSubmatch(g.URL)
	const minMatches = 3
	if len(matches) < minMatches {
		panic(errors.Errorf("invalid gist url: %s", g.URL))
	}
	return fmt.Sprintf(hostURLPattern, matches[1])
}

func (c *Config) getGistConfig(url string) *GistConfig {
	for i := range c.Gists {
		if c.Gists[i].URL == url {
			return &c.Gists[i]
		}
	}
	return nil
}

// AutoDiscoveryConfig TODO: actual gist sample from github.com/lemoony.
func AutoDiscoveryConfig() *Config {
	return &Config{
		Enabled: false,
		Gists: []GistConfig{
			{
				Enabled:                   false,
				URL:                       "gist.github.com/<USERNAME>",
				AuthenticationMethod:      AuthMethodNone,
				IncludeTags:               []string{},
				NameMode:                  SnippetNameModeCombinePreferDescription,
				TitleHeaderEnabled:        true,
				HideTitleInPreview:        true,
				RemoveTagsFromDescription: true,
			},
		},
	}
}
