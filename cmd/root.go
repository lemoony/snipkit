package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/phuslu/log"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/lemoony/snipkit/internal/app"
	"github.com/lemoony/snipkit/internal/cache"
	"github.com/lemoony/snipkit/internal/config"
	"github.com/lemoony/snipkit/internal/managers"
	"github.com/lemoony/snipkit/internal/ui"
	"github.com/lemoony/snipkit/internal/utils/logutil"
	"github.com/lemoony/snipkit/internal/utils/system"
	"github.com/lemoony/snipkit/internal/utils/termutil"
)

type setup struct {
	terminal ui.TUI
	provider managers.Provider
	v        *viper.Viper
	system   *system.System
}

func (s *setup) configService() config.Service {
	return config.NewService(
		config.WithTerminal(s.terminal),
		config.WithViper(s.v),
		config.WithSystem(s.system),
	)
}

var (
	_defaultSystem = system.NewSystem()
	_defaultSetup  = setup{
		terminal: ui.NewTUI(),
		provider: managers.NewBuilder(cache.New(_defaultSystem)),
		v:        viper.GetViper(),
		system:   _defaultSystem,
	}
)

type ctxKey string

var (
	_setupKey         = ctxKey("_setupKey")
	_appKey           = ctxKey("_app")
	_configServiceKey = ctxKey("_cfgService")
)

func getAppFromContext(ctx context.Context) app.App {
	return getAppFromContextWith(ctx, nil, true)
}

func getAppFromContextWithConfigMigrationCheck(ctx context.Context, checkNeedsMigration bool) app.App {
	return getAppFromContextWith(ctx, nil, checkNeedsMigration)
}

func getAppFromContextWith(ctx context.Context, output *os.File, checkNeedsMigration bool) app.App {
	if v := ctx.Value(_appKey); v != nil {
		return v.(app.App)
	}

	s := getSetupFromContext(ctx)

	var tui ui.TUI
	if output != nil {
		tui = ui.NewTUI(ui.WithStdio(termutil.Stdio{
			In:  os.Stdin,
			Out: output,
			Err: os.Stderr,
		}))
	} else {
		tui = s.terminal
	}

	return app.NewApp(
		app.WithTUI(tui),
		app.WithConfigService(s.configService()),
		app.WithProvider(s.provider),
		app.WithCheckNeedsConfigMigration(checkNeedsMigration),
	)
}

func getSetupFromContext(ctx context.Context) setup {
	if v := ctx.Value(_setupKey); v != nil {
		if s, ok := v.(*setup); ok {
			return *s
		}
	}
	return _defaultSetup
}

func getConfigServiceFromContext(ctx context.Context) config.Service {
	if v := ctx.Value(_configServiceKey); v != nil {
		return v.(config.Service)
	}

	s := getSetupFromContext(ctx)
	return s.configService()
}

var (
	cfgFile  string
	logLevel string
)

var rootCmd = &cobra.Command{
	Use:   "snipkit",
	Short: "Use your favorite command line manager directly from the terminal",
	Long:  `Snipkit helps you to execute scripts saved in your favorite snippets manager without even leaving the terminal.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ctx ...context.Context) {
	defer handlePanic()
	setDefaultCommandIfNecessary()
	if len(ctx) > 0 {
		cobra.CheckErr(rootCmd.ExecuteContext(ctx[0]))
	} else {
		cobra.CheckErr(rootCmd.Execute())
	}
}

func SetVersion(v string) {
	rootCmd.Version = v
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().
		StringVarP(&cfgFile, "config", "c", _defaultSetup.system.ConfigPath(), "config file")

	rootCmd.PersistentFlags().StringVarP(
		&logLevel,
		"log-level",
		"l",
		log.PanicLevel.String(),
		fmt.Sprintf("log level used for debugging problems (supported values: %s)", logutil.AllLevelsAsString()))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	_defaultSetup.v.SetConfigFile(cfgFile) // use config file from the flag.
	_defaultSetup.v.AutomaticEnv()         // read in environment variables that match

	configureLogging()
}

func configureLogging() {
	logutil.ConfigureDefaultLogger(_defaultSetup.system)
	logutil.SetDefaultLogLevel(logLevel)
}

func handlePanic() {
	if err := recover(); err != nil {
		fmt.Println(err)

		if e, ok := err.(error); ok {
			log.Error().Err(e).Stack().Msgf("Exited with panic error: %s", e)
		} else {
			log.Error().Msgf("Exited with panic: %s", err)
		}
	}
}

func setDefaultCommandIfNecessary() {
	if c, _, _ := rootCmd.Find(os.Args[1:]); c != rootCmd {
		return
	}

	if defaultCommand, ok := getDefaultCommand(); ok {
		flags := os.Args[1:]
		defaultCommandFields := strings.Fields(defaultCommand)
		os.Args = append([]string{os.Args[0]}, defaultCommandFields...)
		if len(flags) > 0 {
			os.Args = append(os.Args, flags...)
		}
	}
}

func getDefaultCommand() (string, bool) {
	flag.StringVarP(&cfgFile, "config", "c", _defaultSetup.system.ConfigPath(), "config file")
	initConfig()

	if cfg, err := _defaultSetup.configService().LoadConfig(); err == nil {
		return cfg.DefaultRootCommand, true
	}

	return "", false
}
