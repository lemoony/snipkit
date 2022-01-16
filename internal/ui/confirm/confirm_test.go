package confirm

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/utils/termtest"
	"github.com/lemoony/snipkit/internal/utils/termutil"
)

func Test_Confirm(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		expected bool
		send     []string
	}{
		{name: "abort", expected: false, send: termtest.Keys(termtest.KeyEnter)},
		{name: "tab / toggle", expected: true, send: termtest.Keys(termtest.KeyTab, termtest.KeyEnter)},
		{name: "toggle twice", expected: false, send: termtest.Keys(termtest.KeyTab, termtest.KeyTab, termtest.KeyEnter)},
		{name: "y", expected: true, send: []string{"y"}},
		{name: "n", expected: false, send: []string{"n"}},
		{name: "esc", expected: false, send: []string{string(rune(27))}},
		{name: "left", expected: true, send: termtest.Keys(termtest.KeyLeft, termtest.KeyEnter)},
		{name: "right", expected: false, send: termtest.Keys(termtest.KeyRight, termtest.KeyEnter)},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			termtest.RunTerminalTest(t, func(c *termtest.Console) {
				for _, r := range tt.send {
					c.SendString(r)
				}
			}, func(stdio termutil.Stdio) {
				result := Confirm("Are you sure?", "Hello", WithIn(stdio.In), WithOut(stdio.Out), WithFullscreen())
				assert.Equal(t, tt.expected, result)
			})
		})
	}
}

func Test_ConfirmFormatting(t *testing.T) {
	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		c.ExpectString("Hello world")
		c.ExpectString("Are you sure?")
		c.SendString("y")
		c.ExpectString("Yes")
	}, func(stdio termutil.Stdio) {
		header := `Hello world`

		result := Confirm("Are you sure?", header,
			WithIn(stdio.In),
			WithOut(stdio.Out),
			WithSelectionColor("#ff0000"),
		)

		assert.True(t, result)
	})
}

func Test_zeroAwareMin(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{name: "a 1, b 1", a: 1, b: 1, expected: 1},
		{name: "a 2, b 1", a: 2, b: 1, expected: 1},
		{name: "a 1, b 2", a: 1, b: 2, expected: 1},
		{name: "a 0, b 2", a: 0, b: 2, expected: 2},
		{name: "a 2, b 0", a: 2, b: 0, expected: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, zeroAwareMin(tt.a, tt.b))
		})
	}
}
