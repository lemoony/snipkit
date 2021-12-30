package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/adrg/xdg"
	"github.com/phuslu/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lemoony/snippet-kit/internal/cli"
	"github.com/lemoony/snippet-kit/internal/ui"
)

type baseDirectory string

func (d baseDirectory) configPath() string {
	return path.Join(d.path(), "/config.yaml")
}

func (d baseDirectory) path() string {
	return string(d)
}

var terminal = ui.NewTerminal()

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "snipkit",
	Short: "Use your favorite command line manager directly from the terminal",
	Long:  `Snipkit helps you to execute scripts saved in your favorite snippets manager without even leaving the terminal.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cli.LookupAndExecuteSnippet(viper.GetViper(), terminal)
	},
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
	viper.SetConfigFile(cfgFile) // use config file from the flag.
	viper.AutomaticEnv()         // read in environment variables that match
}

func defaultBaseDirectory() (baseDirectory, error) {
	u, err := url.Parse(xdg.ConfigHome)
	if err != nil {
		return "", err
	}
	return baseDirectory(path.Join(u.Path, "snipkit/")), nil
}

func configureLogging() {
	if log.IsTerminal(os.Stderr.Fd()) {
		log.DefaultLogger = log.Logger{
			TimeFormat: "15:04:05",
			Caller:     1,
			Writer: &log.ConsoleWriter{
				ColorOutput:    true,
				QuoteString:    true,
				EndWithMessage: true,
			},
		}
	}

	allLevels := []string{
		log.TraceLevel.String(),
		log.DebugLevel.String(),
		log.InfoLevel.String(),
		log.WarnLevel.String(),
		log.ErrorLevel.String(),
		log.FatalLevel.String(),
		log.PanicLevel.String(),
	}

	var logLevel string
	rootCmd.PersistentFlags().StringVarP(
		&logLevel,
		"log-level",
		"l",
		log.PanicLevel.String(),
		fmt.Sprintf("log level used for debugging problems (supported values: %s)", strings.Join(allLevels, ",")))

	for _, level := range allLevels {
		if logLevel == level {
			log.DefaultLogger.Level = log.ParseLevel(level)
		}
	}
}
