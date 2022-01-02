package errorutil

import (
	"emperror.dev/errors"
)

func NewError(cause error, root error) error {
	return errors.Append(cause, errors.WithStackIf(root))
}
