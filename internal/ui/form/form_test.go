package form

import (
	"testing"

	"github.com/stretchr/testify/assert"

	internalModel "github.com/lemoony/snipkit/internal/model"
	"github.com/lemoony/snipkit/internal/utils/termtest"
	"github.com/lemoony/snipkit/internal/utils/termutil"
)

var testFields = []internalModel.Parameter{
	{Key: "message", Name: "Message", Description: "What to print first"},
	{
		Key:         "application",
		Description: "A second information for the terminal",
		Values: []string{
			"The Romans learned from the Greeks",
			"probably marmelada",
			"by the French name cotignac",
			"option 4",
			"optopn 5",
			"option 6",
			"option 7",
			"option 8",
		},
	},
	{
		Key:          "Statement",
		Description:  "A description",
		DefaultValue: "default value",
	},
}

func Test_ShowForm(t *testing.T) {
	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		c.ExpectString("This snippet requires parameters")

		c.Send("hello")
		c.SendKey(termtest.KeyEnter)

		c.SendKey(termtest.KeyDown)
		c.SendKey(termtest.KeyDown)
		c.SendKey(termtest.KeyDown)
		c.SendKey(termtest.KeyUp)
		c.SendKey(termtest.KeyEnter)

		// just apply the default value for the 3rd parameter
		c.SendKey(termtest.KeyEnter)

		// hit enter
		c.SendKey(termtest.KeyEnter)
	}, func(stdio termutil.Stdio) {
		result, ok := Show(testFields, "ok", WithIn(stdio.In), WithOut(stdio.Out))
		assert.Equal(t, true, ok)
		assert.Len(t, result, 3)
		assert.Equal(t, "hello", result[0])
		assert.Equal(t, "probably marmelada", result[1])
		assert.Equal(t, "default value", result[2])
	})
}

func Test_ShowForm_NextTabAndThenCancel(t *testing.T) {
	termtest.RunTerminalTest(t, func(c *termtest.Console) {
		c.ExpectString("This snippet requires parameters")

		// jump all fields
		for i := 0; i < len(testFields); i++ {
			c.SendKey(termtest.KeyTab)
		}

		c.SendKey(termtest.KeyTab) // jump ok button
		c.SendKey(termtest.KeyTab) // jump abort button

		// jump all fields
		for i := 0; i < len(testFields); i++ {
			c.SendKey(termtest.KeyTab) // jump 1. field
		}

		c.SendKey(termtest.KeyTab)   // jump ok button
		c.SendKey(termtest.KeyEnter) // hit abort button
	}, func(stdio termutil.Stdio) {
		result, ok := Show(testFields, "ok", WithIn(stdio.In), WithOut(stdio.Out))
		assert.False(t, ok)
		assert.Len(t, result, 0)
	})
}
