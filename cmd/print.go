package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Prints the snippet on stdout",
	Long:  `Prints the selected snippet on stdout with all parameters being replaced.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		app := getAppFromContext(cmd.Context())
		snippet := app.LookupAndCreatePrintableSnippet()
		fmt.Println(snippet)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
}
