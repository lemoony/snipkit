package config

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	oldConfigPath = "../../testdata/migrations/config-1-1-0.yaml"
	newConfigPath = "../../testdata/example-config.yaml"
)

func Test_Migrate(t *testing.T) {
	oldCfg, err := ioutil.ReadFile(oldConfigPath)
	assert.NoError(t, err)

	newCfg, err := ioutil.ReadFile(newConfigPath)
	assert.NoError(t, err)

	actualCfg := Migrate(oldCfg)

	assert.YAMLEq(t, string(newCfg), string(actualCfg))
}
