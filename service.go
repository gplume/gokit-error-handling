package main

import (
	"github.com/pkg/errors"
)

// StringService provides operations on strings.
type StringService interface {
	Uppercase(string) (string, error)
	Count(string) int
}

type stringService struct{}

func (stringService) Uppercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	// panic("P A N I C")
	// return strings.ToUpper(s), nil
	err := errors.New("INSTANCE ERROR")
	return "", errors.Wrap(err, "UPPERCASE error")
	// return "", errors.New("UPPERCASE error")
}

func (stringService) Count(s string) int {
	return len(s)
}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")
