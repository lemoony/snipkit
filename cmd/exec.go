package cmd

import (
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a snippet directly from the terminal",
	Long:  `Execute a snippet directly from the terminal. The output of the commands will be visibile in the terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		app := getAppFromContext(cmd.Context())
		app.LookupAndExecuteSnippet()
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
