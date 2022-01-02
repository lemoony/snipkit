package app

import (
	"errors"
	"io/ioutil"
	"path"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snippet-kit/internal/config"
)

func Test_NewApp_NoConfigFile(t *testing.T) {
	v := viper.NewWithOptions()

	_, err := NewApp(v)
	assert.True(t, errors.Is(err, config.ErrNoConfigFound))
}

func Test_NewAppInvalidConfigFile(t *testing.T) {
	cfgFile := path.Join(t.TempDir(), "invalid-config")
	assert.NoError(t, ioutil.WriteFile(cfgFile, []byte("invalid"), 0o600))

	v := viper.NewWithOptions()
	v.SetConfigFile(cfgFile)

	_, err := NewApp(v)
	assert.True(t, errors.Is(err, config.ErrInvalidConfig))
}
