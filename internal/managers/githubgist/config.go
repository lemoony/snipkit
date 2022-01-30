package githubgist

import "github.com/lemoony/snipkit/internal/utils/system"

type AuthMethod string

const (
	AuthMethodNone  = AuthMethod("None")
	AuthMethodToken = AuthMethod("Token")
	AuthMethodOAuth = AuthMethod("OAuth")
)

type Config struct {
	Enabled bool         `yaml:"enabled" head_comment:"If set to false, github gist is disabled completely."`
	Gists   []GistConfig `yaml:"gists" head_comment:"You can define multiple independent Github Gist sources."`
}

type GistConfig struct {
	Enabled              bool       `yaml:"enabled" head_comment:"If set to false, this github gist url is ignored."`
	URL                  string     `yaml:"url" head_comment:"Endpoint url for the gists. Most likely, you want this to point to an user account."`
	AuthenticationMethod AuthMethod `yaml:"authenticationMethod" head_comment:"Supported values: Public, OAuth, Token. Default value: None (which means no authentication). In order to retrieve secret gists, you must be authenticated."`
	IncludeTags          []string   `yaml:"includeTags" head_comment:"If this list is not empty, only those gists that match the listed tags will be provided to you."`
	SuffixRegex          []string   `yaml:"suffixRegex" head_comment:"Only gist files with endings which match one of the listed suffixes will be considered."`
}

func AutoDiscoveryConfig(system *system.System) *Config {
	return &Config{
		Enabled: false,
		Gists: []GistConfig{
			{
				Enabled:              false,
				URL:                  "https://api.github.com/users/<yourUser>/gists",
				AuthenticationMethod: AuthMethodNone,
				IncludeTags:          []string{},
			},
		},
	}
}
