package testdata

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ConfigVersion string

const (
	ConfigV110 = ConfigVersion("1.1.0")
	ConfigV111 = ConfigVersion("1.1.1")
	Latest     = ConfigV111
)

func ConfigPath(t *testing.T, cfgVersion ConfigVersion) string {
	return path.Join(
		absolutePath(t),
		"migrations",
		fmt.Sprintf("config-%s.yaml", strings.ReplaceAll(string(cfgVersion), ".", "-")),
	)
}

func ConfigBytes(t *testing.T, cfgVersion ConfigVersion) []byte {
	fmt.Println(ConfigPath(t, cfgVersion))
	bytes, err := ioutil.ReadFile(ConfigPath(t, cfgVersion))
	assert.NoError(t, err)
	return bytes
}

func absolutePath(t *testing.T) string {
	_, filename, _, ok := runtime.Caller(1)
	fmt.Println(filename)
	assert.True(t, ok)
	return filepath.Dir(filename)
}
