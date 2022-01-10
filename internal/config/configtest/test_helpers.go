package configtest

import (
	"os"
	"path"
	"testing"

	"emperror.dev/errors"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"

	"github.com/lemoony/snippet-kit/internal/config"
	"github.com/lemoony/snippet-kit/internal/providers"
	"github.com/lemoony/snippet-kit/internal/providers/fslibrary"
	"github.com/lemoony/snippet-kit/internal/providers/snippetslab"
	"github.com/lemoony/snippet-kit/internal/ui"
)

const fileModeConfig = os.FileMode(0o600)

// Option configures an App.
type Option interface {
	apply(c *config.Config)
}

// terminalOptionFunc wraps a func so that it satisfies the Option interface.
type optionFunc func(c *config.Config)

func (f optionFunc) apply(c *config.Config) {
	f(c)
}

// WithAdapter allows to modify the config after creation but before marshaling.
func WithAdapter(adapter func(c *config.Config)) Option {
	return optionFunc(func(c *config.Config) {
		adapter(c)
	})
}

func NewTestConfigFilePath(t *testing.T, fs afero.Fs, options ...Option) string {
	t.Helper()

	versionWrapper := NewTestConfig()
	for _, o := range options {
		o.apply(&versionWrapper.Config)
	}

	cfgFilePath := path.Join(t.TempDir(), "temp-config.yaml")

	bytes, err := yaml.Marshal(versionWrapper)
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal config"))
	}

	err = afero.WriteFile(fs, cfgFilePath, bytes, fileModeConfig)
	if err != nil {
		panic(errors.Wrap(err, "failed to write test config file"))
	}

	return cfgFilePath
}

func NewTestConfig() config.VersionWrapper {
	return config.VersionWrapper{
		Version: "1.0",
		Config: config.Config{
			Editor: "test-editor",
			Style:  ui.DefaultConfig(),
			Providers: providers.Config{
				SnippetsLab: snippetslab.Config{
					Enabled: false,
				},
				FsLibrary: fslibrary.Config{
					Enabled: false,
				},
			},
		},
	}
}
