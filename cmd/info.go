package cmd

import (
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Provides useful information about the snipkit configuration",
	Long: `This command is useful to view the current configuration of SnipKit, 
helping to debug any issues you may experience`,
	RunE: func(cmd *cobra.Command, args []string) error {
		app, err := getAppFromContext(cmd.Context())
		if err != nil {
			return err
		}

		return app.Info()
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
