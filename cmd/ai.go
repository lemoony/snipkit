package cmd

import (
	"github.com/spf13/cobra"
)

var aiCmd = &cobra.Command{
	Use:   "ai",
	Short: "Generate a script based on a user prompt (short-hand alias for assistant start).",
	Long:  `Generate a script based on a user prompt and either copy it to the clipboard or execute it directly (short-hand alias for assistant start).`,
	Run: func(cmd *cobra.Command, args []string) {
		app := getAppFromContext(cmd.Context())
		app.CreateSnippetWithAI()
	},
}

func init() {
	rootCmd.AddCommand(aiCmd)
}
