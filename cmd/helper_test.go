package cmd

import (
	"bytes"
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/lemoony/snippet-kit/internal/app"
	"github.com/lemoony/snippet-kit/internal/providers"
	"github.com/lemoony/snippet-kit/internal/ui"
	"github.com/lemoony/snippet-kit/internal/utils"
	"github.com/lemoony/snippet-kit/mocks"
)

type testSetup struct {
	system           *utils.System
	v                *viper.Viper
	providersBuilder providers.Builder
}

func newTestSetup() *testSetup {
	system := utils.NewTestSystem()
	v := viper.New()
	v.SetFs(system.Fs)

	return &testSetup{
		system:           system,
		v:                v,
		providersBuilder: providers.NewBuilder(),
	}
}

// option configures an App.
type option interface {
	apply(t *testSetup)
}

// terminalOptionFunc wraps a func so that it satisfies the option interface.
type optionFunc func(t *testSetup)

func (f optionFunc) apply(t *testSetup) {
	f(t)
}

// withConfigFilePath sets the config file path for the App.
func withConfigFilePath(cfgFilePath string) option {
	return optionFunc(func(t *testSetup) {
		t.v.SetConfigFile(cfgFilePath)
	})
}

// withSystem sets the system for test.
func withSystem(s *utils.System) option {
	return optionFunc(func(t *testSetup) {
		t.system = s
		t.v.SetFs(s.Fs)
	})
}

// withProviders sets the providers for test.
func withProviders(p ...providers.Provider) option {
	return optionFunc(func(t *testSetup) {
		builder := mocks.Builder{}
		builder.On("BuildProvider", mock.Anything, mock.Anything).Return(p, nil)
		t.providersBuilder = &builder
	})
}

func runMockedTest(t *testing.T, args []string, app app.App) error {
	t.Helper()

	ctx := context.WithValue(context.Background(), _appKey, app)

	defer rootCmd.ResetFlags()
	rootCmd.SetArgs(args)

	return rootCmd.ExecuteContext(ctx)
}

func runVT10XCommandTest(
	t *testing.T, args []string, hasError bool, test func(*expect.Console, *setup), options ...option,
) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("Windows does not support psuedoterminals")
	}

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
		system:           testSetup.system,
		v:                testSetup.v,
		providersBuilder: testSetup.providersBuilder,
		terminal:         ui.NewTerminal(ui.WithStdio(terminal.Stdio{In: c.Tty(), Out: c.Tty(), Err: c.Tty()})),
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
