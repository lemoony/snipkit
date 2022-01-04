package cmd

import (
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a snippet directly from the terminal",
	Long:  `Execute a snippet directly from the terminal. The output of the commands will be visibile in the terminal.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		app, err := getAppFromContext(cmd.Context())
		if err != nil {
			return err
		}

		return app.LookupAndExecuteSnippet()
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
