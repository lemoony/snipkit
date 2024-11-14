package configtest

import (
	"os"
	"path"
	"testing"

	"emperror.dev/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"github.com/lemoony/snipkit/internal/assistant"
	"github.com/lemoony/snipkit/internal/assistant/openai"
	"github.com/lemoony/snipkit/internal/config"
	"github.com/lemoony/snipkit/internal/managers"
	"github.com/lemoony/snipkit/internal/managers/fslibrary"
	"github.com/lemoony/snipkit/internal/managers/snippetslab"
	"github.com/lemoony/snipkit/internal/ui"
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

	cfgFilePath := path.Join(t.TempDir(), "config.yaml")

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
		Version: config.Version,
		Config: config.Config{
			Editor: "test-editor",
			Style:  ui.DefaultConfig(),
			Assistant: assistant.Config{
				SaveMode: assistant.SaveModeFsLibrary,
				OpenAI: &openai.Config{
					Enabled: true,
				},
			},
			Manager: managers.Config{
				SnippetsLab: &snippetslab.Config{
					Enabled: false,
				},
				FsLibrary: &fslibrary.Config{
					Enabled: false,
				},
			},
		},
	}
}

func SetSnipkitHomeEnv(t *testing.T, val string) {
	t.Helper()
	assert.NoError(t, os.Setenv("SNIPKIT_HOME", val))
}

func ResetSnipkitHome(t *testing.T) {
	t.Helper()
	assert.NoError(t, os.Unsetenv("SNIPKIT_HOME"))
}
