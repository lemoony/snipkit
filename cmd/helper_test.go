package cmd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/lemoony/snipkit/internal/app"
	"github.com/lemoony/snipkit/internal/config"
	"github.com/lemoony/snipkit/internal/managers"
	"github.com/lemoony/snipkit/internal/ui"
	"github.com/lemoony/snipkit/internal/utils/termtest"
	"github.com/lemoony/snipkit/internal/utils/termutil"
	mocks "github.com/lemoony/snipkit/mocks/managers"
)

type testSetup struct {
	app           app.App
	configService config.Service
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

func runMockedTest(t *testing.T, args []string, options ...option) error {
	t.Helper()

	ts := &testSetup{}
	for _, o := range options {
		o.apply(ts)
	}

	ctx := context.Background()
	if ts.app != nil {
		ctx = context.WithValue(ctx, _appKey, ts.app)
	}
	if ts.configService != nil {
		ctx = context.WithValue(ctx, _configServiceKey, ts.configService)
	}

	defer rootCmd.ResetFlags()
	rootCmd.SetArgs(args)

	// need to re-init the config in case custom flags have been apssed
	rootCmd.ResetFlags()
	res := rootCmd.ExecuteContext(ctx)

	return res
}

func runTerminalTest(t *testing.T, args []string, setup setup, hasError bool, test func(*termtest.Console)) {
	t.Helper()
	termtest.RunTerminalTest(t, test, func(stdio termutil.Stdio) {
		prevIn, prevOut, prevErr := rootCmd.InOrStdin(), rootCmd.OutOrStdout(), rootCmd.ErrOrStderr()
		defer func() {
			rootCmd.SetIn(prevIn)
			rootCmd.SetOut(prevOut)
			rootCmd.SetErr(prevErr)
		}()

		rootCmd.SetIn(stdio.In)
		rootCmd.SetOut(stdio.Out)
		rootCmd.SetErr(stdio.Err)

		s := setup
		s.terminal = ui.NewTUI(ui.WithStdio(termutil.Stdio{In: stdio.In, Out: stdio.Out, Err: stdio.Err}))

		defer rootCmd.ResetFlags()
		rootCmd.SetArgs(args)
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
