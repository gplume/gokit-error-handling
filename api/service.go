package api

import (
	"strings"
	"unicode/utf8"

	"github.com/pkg/errors"
)

// StringSvc provides operations on strings.
type StringSvc interface {
	Uppercase(string) (string, error)
	Count(string) int
}

// StringService ...
type StringService struct{}

// Uppercase ...
func (StringService) Uppercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	// panic("P A N I C")
	return strings.ToUpper(s), nil
	// return "", errors.Wrapf(ErrEmpty, "U P - E R R O R")
}

// Count ...
func (StringService) Count(s string) int {
	return utf8.RuneCountInString(s)
}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")
