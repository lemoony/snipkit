package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/lemoony/snippet-kit/internal/cli"
)

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Prints the snippet on stdout",
	Long:  `Prints the selected snippet on stdout with all parameters being replaced.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		snippet, err := cli.LookupAndCreatePrintableSnippet()
		if err != nil {
			return err
		}
		fmt.Println(snippet)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
}
