package cmd

import (
	"github.com/spf13/cobra"
)

var (
	execCmdPrintFlag   = false
	execCmdConfirmFlag = false
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a snippet directly from the terminal",
	Long:  `Execute a snippet directly from the terminal. The output of the commands will be visibile in the terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		app := getAppFromContext(cmd.Context())
		app.LookupAndExecuteSnippet(execCmdConfirmFlag, execCmdPrintFlag)
	},
}

func init() {
	execCmd.PersistentFlags().BoolVarP(
		&execCmdPrintFlag,
		"print",
		"p",
		false,
		"print the command before execution on stdout",
	)

	execCmd.PersistentFlags().BoolVar(
		&execCmdConfirmFlag,
		"confirm",
		false,
		"the command is printed on stdout before execution for confirmation",
	)

	rootCmd.AddCommand(execCmd)
}
