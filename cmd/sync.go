package cmd

import (
	"github.com/spf13/cobra"
)

// syncCmd is a short alias for managerSyncCommand.
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronizes all snippet managers",
	Long: `Synchronizes all snippet managers. This updates the cache in case a specific manager requires caching of the 
snippets. Alias for manager sync.`,
	Run: managerSyncCommand.Run,
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
