// +build !race

package ui

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/AlecAivazis/survey/v2/terminal"
	expect "github.com/Netflix/go-expect"
	"github.com/stretchr/testify/assert"
)

// For some reason this test fails with `go test --race` when executing
// it on macOS. We exclude this test from race tests by putting it
// into a file with a corresponding build constraint.
func Test_OpenEditor(t *testing.T) {
	runExpectTest(t, func(c *expect.Console) {
		_, _ = c.Send("iHello world\x1b")
		time.Sleep(time.Second)
		_, _ = c.SendLine(":wq!")
	}, func(stdio terminal.Stdio) {
		term := NewTerminal(WithStdio(stdio))

		testFile := path.Join(t.TempDir(), "testfile")
		_, err := os.Create(testFile)
		assert.NoError(t, err)

		term.OpenEditor(testFile, "")
		bytes, err := ioutil.ReadFile(testFile) //nolint:gosec // potential file inclusion via variable
		assert.NoError(t, err)
		assert.Equal(t, "Hello world\n", string(bytes))
	})
}

func Test_OpenEditor_InvalidCommand(t *testing.T) {
	runExpectTest(t, func(c *expect.Console) {
		// nothing to expect since panic will be handled at application root level
	}, func(stdio terminal.Stdio) {
		term := NewTerminal(WithStdio(stdio))

		testFile := path.Join(t.TempDir(), "testfile")
		_, err := os.Create(testFile)
		assert.NoError(t, err)

		assert.Panics(t, func() {
			term.OpenEditor(testFile, "foo-editor")
		})
	})
}