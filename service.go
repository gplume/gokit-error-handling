package main

import (
	"unicode/utf8"

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
	// _, w, e, _ := runtime.Caller(0)
	// panic("P A N I C")
	return "", errors.Wrapf(ErrEmpty, "U P - E R R O R")
	// return "", NewErr(er, "UPPERCASE error", http.StatusOK)
	// return "", errors.New("UPPERCASE error")
}

func (stringService) Count(s string) int {
	return utf8.RuneCountInString(s)
}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")
