package testutil

import (
	"fmt"
	"runtime/debug"
	"testing"

	"emperror.dev/errors"
	"github.com/stretchr/testify/assert"
)

func AssertPanicsWithError(t *testing.T, expected error, f assert.PanicTestFunc, msgAndArgs ...interface{}) error {
	t.Helper()

	funcDidPanic, panicValue, panickedStack := didPanic(f)
	if !funcDidPanic {
		assert.Fail(t,
			fmt.Sprintf("func %#v should panic\n\tPanic value:\t%#v", f, panicValue), msgAndArgs...,
		)
	}

	if err, ok := panicValue.(error); !ok || !errors.Is(err, expected) {
		assert.Fail(t,
			fmt.Sprintf("func %#v should panic with error:\t%#v\n\tPanic value:\t%#v\n\tPanic stack:\t%s",
				f, expected, panicValue, panickedStack,
			),
			msgAndArgs...,
		)
	} else {
		return err
	}

	return nil
}

// didPanic returns true if the function passed to it panics. Otherwise, it returns false.
func didPanic(f assert.PanicTestFunc) (bool, interface{}, string) {
	didPanic := false
	var message interface{}
	var stack string
	func() {
		defer func() {
			if message = recover(); message != nil {
				didPanic = true
				stack = string(debug.Stack())
			}
		}()

		// call the target function
		f()
	}()

	return didPanic, message, stack
}
