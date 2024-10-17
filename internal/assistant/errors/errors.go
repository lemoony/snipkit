package errors

import (
	"emperror.dev/errors"
)

var ErrorNoClientConfiguredOrEnabled = errors.New("No assistant configured or enabled. Try 'snipkit assistant choose'.")
