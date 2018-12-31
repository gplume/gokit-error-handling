package errs

import "github.com/pkg/errors"

var (
	// ErrInternalServer ...
	ErrInternalServer = errors.New("an internal server error occurred please contact the server's administrator")
	// ErrInvalidBody ...
	ErrInvalidBody = errors.New("invalid body")
	// ErrEmpty is returned when an input string is empty.
	ErrEmpty = errors.New("empty string")
)
