package errs

import (
	"errors"
	"net/http"
)

// Level of the error
type Level int

// Error level definition
const (
	Undefined Level = iota // == 0
	Critical               // 1...
	_
	_
	_
	Average
	_
	_
	_
	Low
	end
)

var (
	// ErrInternalServer ...
	ErrInternalServer = errors.New("an internal server error occurred please contact the server's administrator")
	// ErrInvalidBody ...
	ErrInvalidBody = errors.New("invalid body")
	// ErrJustMessage ...
	ErrJustMessage = &Error{
		Message: "Just a message here...",
	}
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
		Level:   Critical,
		// will be automatically filled at runtime:
		// Caller:
		// Stack:
	}
)
