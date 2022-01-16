package cmd

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lemoony/snipkit/internal/app"
	"github.com/lemoony/snipkit/internal/config"
	"github.com/lemoony/snipkit/internal/managers"
	"github.com/lemoony/snipkit/internal/ui"
	"github.com/lemoony/snipkit/internal/utils/system"
	"github.com/lemoony/snipkit/internal/utils/termutil"
	"github.com/lemoony/snipkit/internal/utils/testutil"
	mocks "github.com/lemoony/snipkit/mocks/managers"
)

type testSetup struct {
	system        *system.System
	v             *viper.Viper
	provider      managers.Provider
	app           app.App
	configService config.Service
}

func newTestSetup() *testSetup {
	system := testutil.NewTestSystem()
	v := viper.New()
	v.SetFs(system.Fs)

	return &testSetup{
		system:   system,
		v:        v,
		provider: managers.NewBuilder(),
	}
}

type option interface {
	apply(t *testSetup)
}

type optionFunc func(t *testSetup)

func (f optionFunc) apply(t *testSetup) {
	f(t)
}

func withConfigFilePath(cfgFilePath string) option {
	return optionFunc(func(t *testSetup) {
		t.v.SetConfigFile(cfgFilePath)
	})
}

func withSystem(s *system.System) option {
	return optionFunc(func(t *testSetup) {
		t.system = s
		t.v.SetFs(s.Fs)
	})
}

func withManager(m ...managers.Manager) option {
	return optionFunc(func(t *testSetup) {
		provider := mocks.Provider{}
		provider.On("CreateManager", mock.Anything, mock.Anything).Return(m, nil)
		t.provider = &provider
	})
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

	testSetup := newTestSetup()
	for _, o := range options {
		o.apply(testSetup)
	}

	ctx := context.Background()
	if testSetup.app != nil {
		ctx = context.WithValue(ctx, _appKey, testSetup.app)
	}
	if testSetup.configService != nil {
		ctx = context.WithValue(ctx, _configServiceKey, testSetup.configService)
	}

	defer rootCmd.ResetFlags()
	rootCmd.SetArgs(args)

	// need to re-init the config in case custom flags have been apssed
	rootCmd.ResetFlags()
	res := rootCmd.ExecuteContext(ctx)

	fmt.Println(res)

	return res
}

func runVT10XCommandTest(
	t *testing.T, args []string, hasError bool, test func(*expect.Console, *setup), options ...option,
) {
	t.Helper()

	// Multiplex output to a buffer as well for the raw bytes.
	buf := new(bytes.Buffer)
	c, state, err := vt10x.NewVT10XConsole(
		expect.WithStdout(buf),
		expect.WithDefaultTimeout(time.Second*2),
	)
	require.Nil(t, err)
	defer func() {
		_ = c.Close()
	}()

	donec := make(chan struct{})

	rootCmd.SetIn(c.Tty())
	rootCmd.SetOut(c.Tty())
	rootCmd.SetErr(c.Tty())

	testSetup := newTestSetup()
	for _, o := range options {
		o.apply(testSetup)
	}

	_setup := &setup{
		system:   testSetup.system,
		v:        testSetup.v,
		provider: testSetup.provider,
		terminal: ui.NewTerminal(ui.WithStdio(termutil.Stdio{In: c.Tty(), Out: c.Tty(), Err: c.Tty()})),
	}

	defer rootCmd.ResetFlags()
	rootCmd.SetArgs(args)
	err = rootCmd.ExecuteContext(context.WithValue(context.Background(), _setupKey, _setup))

	go func() {
		defer close(donec)
		test(c, _setup)
	}()

	<-donec
	// Close the slave end of the pty, and read the remaining bytes from the master end.
	assert.NoError(t, c.Tty().Close())

	t.Logf("Raw output: %q", buf.String())

	// Dump the terminal's screen.
	t.Logf("\n%s", expect.StripTrailingEmptyLines(state.String()))

	if hasError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
	}
}
