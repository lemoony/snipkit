package app

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snippet-kit/internal/ui/uimsg"
	"github.com/lemoony/snippet-kit/mocks"
)

func Test_NewApp_NoConfigFile(t *testing.T) {
	v := viper.NewWithOptions()

	term := mocks.Terminal{}
	term.On(mocks.PrintError, uimsg.NoConfig()).Return()

	_, err := NewApp(v, WithTerminal(&term))
	assert.NoError(t, err)
}

func Test_NewAppInvalidConfigFile(t *testing.T) {
	cfgFile := path.Join(t.TempDir(), "invalid-config")
	assert.NoError(t, ioutil.WriteFile(cfgFile, []byte("invalid"), 0o600))

	v := viper.NewWithOptions()
	v.SetConfigFile(cfgFile)

	_, err := NewApp(v)
	assert.Error(t, err)
}
