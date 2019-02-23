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
	Fatal
	Error
	Warn
	Info
	Debug
	Trace
	End
)

// Level or importance of errors, 'Debug', 'Trace' and 'Low'
// are not logged by default
func (l level) String() string {
	switch l {
	case Undefined:
		return "Undefined"
	case Fatal:
		return "Fatal"
	case Error:
		return "Error"
	case Warn:
		return "Warn"
	case Info:
		return "Info"
	case Debug:
		return "Debug"
	case Trace:
		return "Trace"
	default:
		return strconv.Itoa(int(l))
	}
}

// Collection of usual errors for convenience and easy testing.
// TODO: complete this list
// These fields will be automatically filled at runtime:
// Caller:
// Stack:
var (
	// ErrInternalServer the basic 500 throw it all error...
	ErrInternalServer = &Err{
		Message: "an internal server error occurred please contact the server's administrator",
		Code:    http.StatusInternalServerError,
		Level:   Error,
	}

	ErrNotFound = &Err{
		Message: "not found",
		Code:    http.StatusNotFound,
		Level:   Info,
	}

	// ErrInvalidBody is used when the payload is wrong
	ErrInvalidBody = &Err{
		Message: "invalid body",
		Code:    http.StatusBadRequest,
		Level:   Info,
	}

	// ErrInvalidParameter throwed when a query required parameter is missing or wrong
	ErrInvalidParameter = &Err{
		Message: "invalid parameter",
		Code:    http.StatusBadRequest,
		Level:   Info,
	}

	// ErrEmpty is returned when an input string is empty.
	ErrEmptyParam = &Err{
		Message: "empty parameter",
		Code:    http.StatusBadRequest,
		Err:     errors.New("empty parameter"),
		Level:   Info,
	}
)
