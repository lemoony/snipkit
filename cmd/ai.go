package cmd

import (
	"github.com/spf13/cobra"
)

var aiCmd = &cobra.Command{
	Use:     "ai",
	Short:   "Generate a script by means of AI",
	Long:    `Generate a script by means of AI and either copy it to the clipboard or execute it directly.`,
	Aliases: []string{"gen"},
	Run: func(cmd *cobra.Command, args []string) {
		app := getAppFromContext(cmd.Context())
		app.CreateSnippetWithAI()
	},
}

func init() {
	rootCmd.AddCommand(aiCmd)
}
