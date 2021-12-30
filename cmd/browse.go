package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lemoony/snippet-kit/internal/cli"
)

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Browse all snippets without executing them",
	Long: `Browse all available snippets without executing them after pressing enter. This is a way to explore your library
in a safe way in case executing some scripts by accident would have undesirable effects.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := cli.LookupSnippet(viper.GetViper(), terminal)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(browseCmd)
}
