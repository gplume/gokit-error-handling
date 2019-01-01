package errs

import (
	"errors"
	"net/http"
)

var (
	// ErrInternalServer ...
	ErrInternalServer = errors.New("an internal server error occurred please contact the server's administrator")
	// ErrInvalidBody ...
	ErrInvalidBody = errors.New("invalid body")
	// ErrEmpty is returned when an input string is empty.
	ErrEmpty = &Error{
		Err:     errors.New("empty parameter"),
		Code:    http.StatusBadRequest,
		Message: "invalid body",
	}
	// ErrSpecific is thrown in case of specific error
	ErrSpecific = &Error{
		Message: "Message for the specific error",
		Code:    http.StatusBadRequest,
		// will be filled at runtime but still can be here:
		Err: errors.New("s p e c i f i c"),
		// automatic from here:
		// Caller:
		// Stack:
	}
)
