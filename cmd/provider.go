package cmd

import (
	"github.com/spf13/cobra"
)

var providerCmd = &cobra.Command{
	Use:   "provider",
	Short: "Manage providers, the snippets manager snipkit connects to",
}

var providerAddCommand = &cobra.Command{
	Use:   "add",
	Short: "Add a new snippet manager provider",
	Long: `Add a new snippet manager provider to your config. SnipKit will connect to it and provide all snippets to you
which meet certain criteria.`,
	Run: func(cmd *cobra.Command, args []string) {
		getAppFromContext(cmd.Context()).AddProvider()
	},
}

func init() {
	rootCmd.AddCommand(providerCmd)

	providerCmd.AddCommand(providerAddCommand)
}
