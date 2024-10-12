package cmd

import (
	"github.com/spf13/cobra"
)

var assistantCmd = &cobra.Command{
	Use:     "ai",
	Short:   "Generate a script by means of AI",
	Long:    `Generate a script by means of AI and either copy it to the clipboard or execute it directly.`,
	Aliases: []string{"assistant", "gen"},
	Run: func(cmd *cobra.Command, args []string) {
		app := getAppFromContext(cmd.Context())
		app.CreateSnippetWithAI()
	},
}

func init() {
	rootCmd.AddCommand(assistantCmd)
}
