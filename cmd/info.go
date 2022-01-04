package cmd

import (
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Provides useful information about the snipkit configuration",
	Long: `This command is useful to view the current configuration of SnipKit, 
helping to debug any issues you may experience`,
	Run: func(cmd *cobra.Command, args []string) {
		getAppFromContext(cmd.Context()).Info()
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
