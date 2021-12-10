package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lemoony/snippet-kit/internal/cli"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage your snipkit configuration file",
}

var configInitCommand = &cobra.Command{
	Use:   "init",
	Short: "Initializes the snipkit config",
	Long: `A snipkit configuration file will be generated at a default location best suited for your operation system.
Snipkit will try to detect any supported snippet manager application and configure them accordingly.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cli.ConfigInit(viper.GetViper())
	},
}

var configCleanCommand = &cobra.Command{
	Use:   "clean",
	Short: "Deletes the snipkit config",
	Long: `The snipkit configuration file will be deleted. You have to initialize a new configuration before using snipkit again.
This command is helpful if your configuration file is corrupted or you want to prepare the uninstalling snipkit.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cli.ConfigClean(viper.GetViper())
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configInitCommand)
	configCmd.AddCommand(configCleanCommand)
}
