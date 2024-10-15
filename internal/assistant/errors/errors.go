package errors

import (
	"emperror.dev/errors"
)

var (
	ErrorNoClientConfiguredOrEnabled = errors.New("No assistant configured or enabled")
	ErrorNoOrInvalidAPIKey           = errors.New("No or invalid api key")
)
