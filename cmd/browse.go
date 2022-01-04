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
		app, err := getAppFromContext(cmd.Context())
		if err != nil {
			return err
		}

		_, err = app.LookupSnippet()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(browseCmd)
}
