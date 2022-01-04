package cmd

import (
	"github.com/spf13/cobra"
)

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Browse all snippets without executing them",
	Long: `Browse all available snippets without executing them after pressing enter. This is a way to explore your library
in a safe way in case executing some scripts by accident would have undesirable effects.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		app := getAppFromContext(cmd.Context())
		_ = app.LookupSnippet()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(browseCmd)
}
