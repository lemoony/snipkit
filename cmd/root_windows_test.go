package cmd

func runCommandTest(t *testing.T, args []string, hasError bool, test func(*expect.Console)) {
	t.Helper()
	t.Skip("Windows does not support psuedoterminals")
}
