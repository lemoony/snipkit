package cmd

import (
	"github.com/spf13/cobra"
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
	Run: func(cmd *cobra.Command, args []string) {
		getConfigServiceFromContext(cmd.Context()).Create()
	},
}

var configCleanCommand = &cobra.Command{
	Use:   "clean",
	Short: "Deletes the snipkit config",
	Long: `The snipkit configuration file will be deleted. You have to initialize a new configuration before using snipkit again.
This command is helpful if your configuration file is corrupted or you want to prepare the uninstalling snipkit.`,
	Run: func(cmd *cobra.Command, args []string) {
		getConfigServiceFromContext(cmd.Context()).Clean()
	},
}

var configEditCommand = &cobra.Command{
	Use:   "edit",
	Short: "Edit the snipkit config",
	Long: `The snipkit configuration file will opened in your preferred editor. 
The editor is defined by the $VISUAL or $EDITOR environment variables. Alternatively, 
the editor can also be defined via the snipkit config file. If neither of those are 
present, vim is used.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		getConfigServiceFromContext(cmd.Context()).Edit()
		return nil
	},
}

var configMigrateCommand = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the snipkit config",
	Long: `The snipkit configuration file will be migrated to the latest version.
This command will do nothing if the version of the config file is already up-to-date.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		getConfigServiceFromContext(cmd.Context()).Migrate()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configInitCommand)
	configCmd.AddCommand(configEditCommand)
	configCmd.AddCommand(configMigrateCommand)
	configCmd.AddCommand(configCleanCommand)
}
