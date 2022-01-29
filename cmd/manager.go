package cmd

import (
	"github.com/spf13/cobra"
)

var managerCmd = &cobra.Command{
	Use:   "manager",
	Short: "Manage the snippet managers snipkit connects to",
}

var managerAddCommand = &cobra.Command{
	Use:   "add",
	Short: "Add a new snippet manager",
	Long: `Add a new snippet manager to your config. SnipKit will connect to it and provide all snippets to you
which meet certain criteria.`,
	Run: func(cmd *cobra.Command, args []string) {
		getAppFromContext(cmd.Context()).AddManager()
	},
}

func init() {
	rootCmd.AddCommand(managerCmd)

	managerCmd.AddCommand(managerAddCommand)
}
