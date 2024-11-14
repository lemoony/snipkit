//go:build demo

package cmd

func init() {
	generateCmd.PersistentFlags().StringArrayVar(
		&assistantDemoScriptFlag,
		"demo-script",
		[]string{},
		"Path to a fixed script (demo purposes only)",
	)

	generateCmd.PersistentFlags().IntVar(
		&assistantDemoWaitFlag,
		"demo-wait-seconds",
		0,
		"Seconds to wait before showing the script (demo purposes only)",
	)
}
