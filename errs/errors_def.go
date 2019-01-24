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
	Critical
	High
	Medium
	Low
	UserOnly
	End
)

// Level or importance of errors, 'UserOnly' and 'Low'
// are not logged by default
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
	case UserOnly:
		return "UserOnly"
	case Low:
		return "Low"
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
	ErrInternalServer = &Error{
		Message: "an internal server error occurred please contact the server's administrator",
		Code:    http.StatusInternalServerError,
		Level:   Critical,
	}

	ErrNotFound = &Error{
		Message: "not found",
		Code:    http.StatusNotFound,
		Level:   Medium,
	}

	// ErrInvalidBody is used when the payload is wrong
	ErrInvalidBody = &Error{
		Message: "invalid body",
		Code:    http.StatusBadRequest,
		Level:   UserOnly,
	}

	// ErrInvalidParameter throwed when a query required parameter is missing or wrong
	ErrInvalidParameter = &Error{
		Message: "invalid parameter",
		Code:    http.StatusBadRequest,
		Level:   UserOnly,
	}

	// ErrEmpty is returned when an input string is empty.
	ErrEmptyParam = &Error{
		Message: "empty parameter",
		Code:    http.StatusBadRequest,
		Err:     errors.New("empty parameter"),
		Level:   UserOnly,
	}
)
