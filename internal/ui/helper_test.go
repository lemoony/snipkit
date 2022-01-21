package ui

import (
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
)

func runScreenTest(t *testing.T, procedure func(s tcell.Screen), test func(s tcell.SimulationScreen)) {
	t.Helper()
	screen := tcell.NewSimulationScreen("")
	assert.NoError(t, screen.Init())

	donec := make(chan struct{})
	go func() {
		defer close(donec)
		time.Sleep(time.Millisecond * 50)
		test(screen)
	}()

	procedure(screen)
	<-donec
}
