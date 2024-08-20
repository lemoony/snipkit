package cmd

import (
	"context"
	"os"
	"runtime"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snipkit/internal/app"
	"github.com/lemoony/snipkit/internal/config"
	"github.com/lemoony/snipkit/internal/managers"
	"github.com/lemoony/snipkit/internal/ui"
	"github.com/lemoony/snipkit/internal/utils/system"
	"github.com/lemoony/snipkit/internal/utils/termtest"
	"github.com/lemoony/snipkit/internal/utils/termutil"
	mocks "github.com/lemoony/snipkit/mocks/managers"
)

type testSetup struct {
	app           app.App
	configService config.Service
	v             *viper.Viper
	system        *system.System
}

type option interface {
	apply(t *testSetup)
}

type optionFunc func(t *testSetup)

func (f optionFunc) apply(t *testSetup) {
	f(t)
}

func withConfigService(configService config.Service) option {
	return optionFunc(func(t *testSetup) {
		t.configService = configService
	})
}

func withApp(app app.App) option {
	return optionFunc(func(t *testSetup) {
		t.app = app
	})
}

func withViper(v *viper.Viper) option {
	return optionFunc(func(t *testSetup) {
		t.v = v
	})
}

func withSystem(s *system.System) option {
	return optionFunc(func(t *testSetup) {
		t.system = s
	})
}

func runExecuteTest(t *testing.T, args []string, options ...option) {
	t.Helper()

	ts := &testSetup{}
	for _, o := range options {
		o.apply(ts)
	}

	rootCmd.SetArgs(nil)
	rootCmd.SetContext(context.TODO())

	_preSetup := _defaultSetup
	_preCfgFile := cfgFile
	defer func() {
		_defaultSetup = _preSetup
		cfgFile = _preCfgFile
	}()

	_defaultSetup = setup{
		v:        _defaultSetup.v,
		system:   _defaultSetup.system,
		provider: _defaultSetup.provider,
		terminal: _defaultSetup.terminal,
	}

	ctx := context.Background()

	if ts.app != nil {
		ctx = context.WithValue(ctx, _appKey, ts.app)
	}
	if ts.configService != nil {
		ctx = context.WithValue(ctx, _configServiceKey, ts.configService)
	}
	if ts.v != nil {
		_defaultSetup.v = ts.v
		cfgFile = ts.v.ConfigFileUsed()
	}
	if ts.system != nil {
		_defaultSetup.system = ts.system
	}

	_, filename, _, _ := runtime.Caller(1)
	os.Args = append([]string{filename}, args...)

	Execute(ctx)
}

func runTerminalTest(t *testing.T, args []string, setup setup, hasError bool, test func(*termtest.Console)) {
	t.Helper()
	termtest.RunTerminalTest(t, test, func(stdio termutil.Stdio) {
		prevIn, prevOut, prevErr := rootCmd.InOrStdin(), rootCmd.OutOrStdout(), rootCmd.ErrOrStderr()
		defer func() {
			rootCmd.ResetFlags()
			rootCmd.SetIn(prevIn)
			rootCmd.SetOut(prevOut)
			rootCmd.SetErr(prevErr)

			for i := range rootCmd.Commands() {
				rootCmd.Commands()[i].ResetFlags()
				//nolint:staticcheck // required for testing since cobra will use an old context instead
				rootCmd.Commands()[i].SetContext(nil)
			}
		}()

		rootCmd.SetIn(stdio.In)
		rootCmd.SetOut(stdio.Out)
		rootCmd.SetErr(stdio.Err)
		rootCmd.SetArgs(args)

		s := setup
		s.terminal = ui.NewTUI(ui.WithStdio(termutil.Stdio{In: stdio.In, Out: stdio.Out, Err: stdio.Err}))

		err := rootCmd.ExecuteContext(context.WithValue(context.Background(), _setupKey, &s))

		if hasError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	})
}

func testProviderForManager(manager managers.Manager) managers.Provider {
	provider := mocks.Provider{}
	provider.On("CreateManager", mock.Anything, mock.Anything).Return([]managers.Manager{manager}, nil)
	return &provider
}
