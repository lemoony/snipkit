package cmd

import (
	"github.com/spf13/cobra"
)

var aiCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generate a script based on a user prompt (short-hand alias for assistant generate).",
	Long:    `Generate a script based on a user prompt and either copy it to the clipboard or execute it directly (short-hand alias for assistant generate).`,
	Aliases: []string{"ai"},
	Run: func(cmd *cobra.Command, args []string) {
		getAppFromContext(cmd.Context()).GenerateSnippetWithAssistant()
	},
}

func init() {
	rootCmd.AddCommand(aiCmd)
}
