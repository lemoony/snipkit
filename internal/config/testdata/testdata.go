package testdata

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ConfigVersion string

const (
	ConfigV100          = ConfigVersion("1.0.0")
	ConfigV110          = ConfigVersion("1.1.0")
	ConfigV111          = ConfigVersion("1.1.1")
	ConfigV111FsLibrary = ConfigVersion("1.1.1-fslibrary")
	ConfigV120          = ConfigVersion("1.2.0")
	ConfigV120FsLibrary = ConfigVersion("1.2.0-fslibrary")
	ConfigV120Providers = ConfigVersion("1.2.0-providers")
	ConfigV130          = ConfigVersion("1.3.0")
	ConfigV130Providers = ConfigVersion("1.3.0-providers")
	Latest              = ConfigV130

	Example = ConfigVersion("example-config.yaml")
)

func ConfigPath(t *testing.T, cfgVersion ConfigVersion) string {
	if cfgVersion == Example {
		return path.Join(absolutePath(t), "example-config.yaml")
	}

	return path.Join(
		absolutePath(t),
		"migrations",
		fmt.Sprintf("config-%s.yaml", strings.ReplaceAll(string(cfgVersion), ".", "-")),
	)
}

func ConfigBytes(t *testing.T, cfgVersion ConfigVersion) []byte {
	bytes, err := os.ReadFile(ConfigPath(t, cfgVersion))
	assert.NoError(t, err)
	return bytes
}

func absolutePath(t *testing.T) string {
	_, filename, _, ok := runtime.Caller(1)
	assert.True(t, ok)
	return filepath.Dir(filename)
}
