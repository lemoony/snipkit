package assertutil

import (
	"fmt"
	"runtime/debug"
	"testing"

	"emperror.dev/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/lemoony/snipkit/internal/model"
)

func AssertExists(t *testing.T, fs afero.Fs, path string, expected bool) {
	t.Helper()
	exists, err := afero.Exists(fs, path)
	assert.NoError(t, err)
	assert.Equal(t, expected, exists)
}

func AssertSnippetsEqual(t *testing.T, expected []model.Snippet, actual []model.Snippet) {
	t.Helper()

	assert.Len(t, actual, len(expected))

	for i, e := range expected {
		a := actual[i]

		assert.Equal(t, e.GetID(), a.GetID())
		assert.Equal(t, e.GetTags(), a.GetTags())
		assert.Equal(t, e.GetLanguage(), a.GetLanguage())
		assert.Equal(t, e.GetTitle(), a.GetTitle())
		assert.Equal(t, e.GetContent(), a.GetContent())
	}
}

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
	panicOccurred := false
	var message interface{}
	var stack string
	func() {
		defer func() {
			if message = recover(); message != nil {
				panicOccurred = true
				stack = string(debug.Stack())
			}
		}()

		// call the target function
		f()
	}()

	return panicOccurred, message, stack
}
