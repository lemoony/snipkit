package cmd

import (
	"github.com/spf13/cobra"
)

var assistantCmd = &cobra.Command{
	Use:   "assistant",
	Short: "SnipKit assistant helps to write executable CLI snippets.",
	Long:  `SnipKit assistant generates a script by means of AI and allows to execute it directly.`,
}

var startCmd = &cobra.Command{
	Use:     "start",
	Short:   "Generate a script based on a user prompt.",
	Long:    `Generate a script based on a user prompt and either copy it to the clipboard or execute it directly.`,
	Aliases: []string{"ai"},
	Run: func(cmd *cobra.Command, args []string) {
		app := getAppFromContext(cmd.Context())
		app.CreateSnippetWithAI()
	},
}

var enableCmd = &cobra.Command{
	Use:     "enable",
	Short:   "Enables a specific assistant provider.",
	Long:    `Enables a specific assistant provider (LLM) by modifying the SnipKit config.`,
	Aliases: []string{"ai"},
	Run: func(cmd *cobra.Command, args []string) {
		getAppFromContext(cmd.Context()).EnableAssistant()
	},
}

func init() {
	rootCmd.AddCommand(assistantCmd)
	assistantCmd.AddCommand(startCmd)
	assistantCmd.AddCommand(enableCmd)
}
