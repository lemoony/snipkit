package cmd

import (
	"context"
	"fmt"
	"net/url"
	"path"

	"github.com/adrg/xdg"
	"github.com/phuslu/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lemoony/snippet-kit/internal/app"
	"github.com/lemoony/snippet-kit/internal/config"
	"github.com/lemoony/snippet-kit/internal/providers"
	"github.com/lemoony/snippet-kit/internal/ui"
	"github.com/lemoony/snippet-kit/internal/utils"
	"github.com/lemoony/snippet-kit/internal/utils/logutil"
)

type setup struct {
	terminal         ui.Terminal
	providersBuilder providers.Builder
	v                *viper.Viper
	system           *utils.System
}

func (s *setup) configService() config.Service {
	return config.NewService(
		config.WithTerminal(s.terminal),
		config.WithViper(s.v),
		config.WithSystem(s.system),
	)
}

var _defaultSetup = setup{
	terminal:         ui.NewTerminal(),
	providersBuilder: providers.NewBuilder(),
	v:                viper.GetViper(),
	system:           utils.NewTestSystem(),
}

type ctxKey string

var (
	_setupKey         = ctxKey("_setupKey")
	_appKey           = ctxKey("_app")
	_configServiceKey = ctxKey("_cfgService")
)

func getAppFromContext(ctx context.Context) app.App {
	if v := ctx.Value(_appKey); v != nil {
		return v.(app.App)
	}

	s := getSetupFromContext(ctx)
	return app.NewApp(
		app.WithTerminal(s.terminal),
		app.WithConfigService(s.configService()),
		app.WithProvidersBuilder(s.providersBuilder),
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

type baseDirectory string

func (d baseDirectory) configPath() string {
	return path.Join(d.path(), "/config.yaml")
}

func (d baseDirectory) path() string {
	return string(d)
}

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "snipkit",
	Short: "Use your favorite command line manager directly from the terminal",
	Long:  `Snipkit helps you to execute scripts saved in your favorite snippets manager without even leaving the terminal.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	if defaultBaseDirectoryPath, err := defaultBaseDirectory(); err != nil {
		cobra.CheckErr(err)
	} else {
		rootCmd.PersistentFlags().
			StringVarP(&cfgFile, "config", "c", defaultBaseDirectoryPath.configPath(), "config file")
	}

	configureLogging()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	_defaultSetup.v.SetConfigFile(cfgFile) // use config file from the flag.
	_defaultSetup.v.AutomaticEnv()         // read in environment variables that match
}

func defaultBaseDirectory() (baseDirectory, error) {
	u, err := url.Parse(xdg.ConfigHome)
	if err != nil {
		return "", err
	}
	return baseDirectory(path.Join(u.Path, "snipkit/")), nil
}

func configureLogging() {
	var logLevel string
	rootCmd.PersistentFlags().StringVarP(
		&logLevel,
		"log-level",
		"l",
		log.PanicLevel.String(),
		fmt.Sprintf("log level used for debugging problems (supported values: %s)", logutil.AllLevelsAsString()))

	logutil.ConfigureDefaultLogger()
	logutil.SetDefaultLogLevel(logLevel)
}
