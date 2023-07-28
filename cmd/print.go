package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	execCmdOutputFlag string
)

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Prints the snippet on stdout",
	Long:  `Prints the selected snippet on stdout with all parameters being replaced.`,
	Run: func(cmd *cobra.Command, args []string) {
		app := getAppFromContext(cmd.Context())
		if snippet, ok := app.LookupAndCreatePrintableSnippet(); ok {
			if execCmdOutputFlag != "" {
				os.WriteFile(execCmdOutputFlag, []byte(snippet), 0644)
			} else {
				fmt.Println(snippet)
			}
		}
	},
}

func init() {
	printCmd.PersistentFlags().StringVarP(
		&execCmdOutputFlag,
		"output",
		"o",
		"",
		fmt.Sprintf("Write snippet to file instead of stdout"))
	rootCmd.AddCommand(printCmd)
}
