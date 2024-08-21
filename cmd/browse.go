package cmd

import (
	"github.com/spf13/cobra"
)

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Browse all snippets without executing them",
	Long: `Browse all available snippets without executing them after pressing enter. This is a way to explore your library
in a safe way in case executing some scripts by accident would have undesirable effects.`,
	Run: func(cmd *cobra.Command, args []string) {
		app := getAppFromContext(cmd.Context())
		_, _ = app.LookupSnippet()
	},
}

func init() {
	rootCmd.AddCommand(browseCmd)
}
