package ui

import (
	"testing"

	"github.com/AlecAivazis/survey/v2/terminal"
	expect "github.com/Netflix/go-expect"
)

func runTest(t *testing.T, procedure func(*expect.Console), test func(terminal.Stdio)) {
	t.Helper()
	t.Skip("Windows does not support psuedoterminals")
}
