package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Prints the snippet on stdout",
	Long:  `Prints the selected snippet on stdout with all parameters being replaced.`,
	Run: func(cmd *cobra.Command, args []string) {
		app := getAppFromContext(cmd.Context())
		if snippet, ok := app.LookupAndCreatePrintableSnippet(); ok {
			fmt.Println(snippet)
		}
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
}
