package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lemoony/snippet-kit/internal/cli"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a snippet directly from the terminal",
	Long:  `Execute a snippet directly from the terminal. The output of the commands will be visibile in the terminal.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cli.LookupAndExecuteSnippet(viper.GetViper(), terminal)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
