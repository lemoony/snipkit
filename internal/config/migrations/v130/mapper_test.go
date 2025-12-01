package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/config/testdata"
)

func Test_Migrate(t *testing.T) {
	oldCfg := testdata.ConfigBytes(t, testdata.ConfigV120)
	newCfg := testdata.ConfigBytes(t, testdata.ConfigV130)
	actualCfg := Migrate(oldCfg)
	assert.YAMLEq(t, string(newCfg), string(actualCfg))
}

func Test_Migrate_WithProviders(t *testing.T) {
	oldCfg := testdata.ConfigBytes(t, testdata.ConfigV120Providers)
	newCfg := testdata.ConfigBytes(t, testdata.ConfigV130Providers)
	actualCfg := Migrate(oldCfg)
	assert.YAMLEq(t, string(newCfg), string(actualCfg))
}

func Test_Migrate_invalidYamlPanic(t *testing.T) {
	assert.Panics(t, func() {
		Migrate([]byte("{"))
	})
}

func Test_Migrate_invalidConfigVersion(t *testing.T) {
	assert.Panics(t, func() {
		Migrate([]byte("version: 3.0.0"))
	})
}
