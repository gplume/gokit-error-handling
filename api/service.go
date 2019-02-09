package api

import (
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/pkg/errors"

	"github.com/gplume/gokit-error-handling/errs"
)

// Service provides operations on strings.
type Service interface {
	Uppercase(string) (string, error)
	Count(string) (int, error)
}

type stringService struct {
}

// NewStringService returns a wrapped service with validation layer
func NewStringService() (Service, error) {
	return newValidationService(Service(stringService{})), nil
}

// Uppercase ...
func (stringService) Uppercase(s string) (string, error) {
	if _, err := strconv.Atoi(s); err == nil {
		return s, errs.New(http.StatusBadRequest, "uppercase some numbers dude, really??", errs.Info)
	}
	if s == "empty" {
		return s, errs.New(errs.ErrEmptyParam)
	}
	if s == "compositing" {
		_, specErr := strconv.Atoi(s)
		return s, errs.New(specErr, errs.ErrInvalidParameter.Message, http.StatusNotAcceptable, errs.Error) // errs.Error overrides defined ErrEmptyParam.Level
	}
	if s == "default" {
		_, err := strconv.Atoi(s)
		return s, errors.WithStack(err)
	}
	if s == "needs" {
		_, err := strconv.Atoi(s)
		return s, errs.New(err)
	}
	if s == "notfound" {
		return s, errs.New(errs.ErrNotFound)
	}
	return strings.ToUpper(s), nil
}

// Count ...
func (stringService) Count(s string) (int, error) {
	return utf8.RuneCountInString(s), nil
}
