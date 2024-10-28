//go:build demo

package cmd

func init() {
	generateCmd.PersistentFlags().StringVar(
		&assistantDemoScriptFlag,
		"demo-script",
		"",
		"Path to a fixed script (demo purposes only)",
	)

	generateCmd.PersistentFlags().IntVar(
		&assistantDemoWaitFlag,
		"demo-wait-seconds",
		0,
		"Seconds to wait before showing the script (demo purposes only)",
	)
}
