package cmd

import (
	"emperror.dev/errors"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

var copyCmd = &cobra.Command{
	Use:     "copy",
	Aliases: []string{"cp"},
	Short:   "Copies the snippet to the clipboard",
	Long:    `Copies the selected snippet to the clipboard for manual execution.`,
	Run: func(cmd *cobra.Command, args []string) {
		app := getAppFromContext(cmd.Context())
		if snippet, ok := app.LookupAndCreatePrintableSnippet(); ok {
			copyToClipboard(snippet)
		}
	},
}

func copyToClipboard(snippet string) {
	if err := clipboard.WriteAll(snippet); err != nil {
		panic(errors.Wrap(errors.WithStack(err), "failed to write to clipboard"))
	}
}

func init() {
	rootCmd.AddCommand(copyCmd)
}
