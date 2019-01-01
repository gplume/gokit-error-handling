package errs

import (
	"errors"
	"net/http"
	"strconv"
)

// level of the error
type level int

// Error level definition
const (
	Undefined level = iota // == 0
	Critical               // 1...
	_
	High // 3
	_
	Medium // 5
	_
	_
	_
	Low // 9
	end
)

func (l level) String() string {
	switch l {
	case Undefined:
		return "Undefined"
	case Critical:
		return "Critical"
	case High:
		return "High"
	case Medium:
		return "Medium"
	case Low:
		return "Low"
	default:
		return strconv.Itoa(int(l))
	}
}

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
