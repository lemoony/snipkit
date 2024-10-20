package cmd

import (
	"github.com/spf13/cobra"
)

var assistantCmd = &cobra.Command{
	Use:   "assistant",
	Short: "SnipKit assistant helps to write executable CLI snippets.",
	Long:  `SnipKit assistant generates a script by means of AI and allows to execute it directly.`,
}

var generateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generate a script based on a user prompt.",
	Long:    `Generate a script based on a user prompt and either copy it to the clipboard or execute it directly.`,
	Aliases: []string{"ai", "create"},
	Run: func(cmd *cobra.Command, args []string) {
		getAppFromContext(cmd.Context()).GenerateSnippetWithAssistant()
	},
}

var choose = &cobra.Command{
	Use:     "choose",
	Short:   "Choose a specific assistant provider.",
	Long:    `Enables a specific assistant provider (LLM) by modifying the SnipKit config.`,
	Aliases: []string{"switch", "enable"},
	Run: func(cmd *cobra.Command, args []string) {
		getAppFromContext(cmd.Context()).EnableAssistant()
	},
}

func init() {
	rootCmd.AddCommand(assistantCmd)
	rootCmd.AddCommand(generateCmd)
	assistantCmd.AddCommand(generateCmd)
	assistantCmd.AddCommand(choose)
}
