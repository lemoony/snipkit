package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
)

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Prints the snippet on stdout",
	Long:  `Prints the selected snippet on stdout with all parameters being replaced.`,
	Run: func(cmd *cobra.Command, args []string) {
		// workaround for subshells: use stderr for output
		// https://github.com/charmbracelet/bubbletea/issues/206
		lipgloss.SetColorProfile(termenv.NewOutput(os.Stderr).Profile)
		app := getAppFromContextWith(cmd.Context(), os.Stderr, true)

		if snippet, ok := app.LookupAndCreatePrintableSnippet(); ok {
			_, _ = fmt.Fprintln(os.Stdout, snippet)
		}
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
}
