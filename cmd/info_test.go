package cmd

import (
	"testing"

	"github.com/Netflix/go-expect"
)

func Test_Info(t *testing.T) {
	runCommandTest(t, []string{"info"}, false, func(console *expect.Console) {
		// TODO
	})
}
