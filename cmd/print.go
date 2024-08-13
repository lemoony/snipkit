package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
)

var (
	printCmdCopyFlag       bool
	printCmdIDFlag         string
	printCmdParametersFlag []string
)

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Prints the snippet on stdout",
	Long:  `Prints the selected snippet on stdout with all parameters being replaced.`,
	Run: func(cmd *cobra.Command, args []string) {
		lipgloss.SetColorProfile(termenv.NewOutput(os.Stderr).Profile)
		app := getAppFromContextWith(cmd.Context(), os.Stderr, true)

		if printCmdIDFlag != "" {
			if snippet, ok := app.FindSnippetAndPrint(printCmdIDFlag, toParameterValues(printCmdParametersFlag)); ok {
				fmt.Println(snippet)
				if printCmdCopyFlag {
					copyToClipboard(snippet)
				}
			}
		} else if snippet, ok := app.LookupAndCreatePrintableSnippet(); ok {
			fmt.Println(snippet)
			if printCmdCopyFlag {
				copyToClipboard(snippet)
			}
		}
	},
}

func init() {
	printCmd.PersistentFlags().BoolVar(
		&printCmdCopyFlag,
		"copy",
		false,
		"copies the snippet to the clipboard",
	)

	printCmd.PersistentFlags().StringVar(
		&printCmdIDFlag,
		"id",
		"",
		"ID of the snippet to print",
	)

	printCmd.PersistentFlags().StringArrayVarP(
		&printCmdParametersFlag,
		"param",
		"p",
		[]string{},
		"Parameter values to be passed to the snippet",
	)

	rootCmd.AddCommand(printCmd)
}
