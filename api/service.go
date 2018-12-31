package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gplume/gokit-error-handling/errs"
)

// Logger is the interface that wraps the basic Log method
type Logger interface {
	Log(keyvals ...interface{}) error
}

// Service provides operations on strings.
type Service interface {
	Uppercase(string) (string, error)
	Count(string) (int, error)
}

// StringService ...
type stringService struct {
	logger Logger
}

// NewStringService returns a wrapped service with validation layer
func NewStringService(logger Logger) (Service, error) {
	svc, err := newStringService(logger)
	if err != nil {
		return nil, err
	}
	svc, err = newValidationService(svc, logger)
	if err != nil {
		return nil, err
	}
	return svc, nil
}

func newStringService(logger Logger) (Service, error) {
	if logger == nil {
		return nil, errors.New("cannot create new string service, logger cannot be nil")
	}
	return stringService{
		logger: logger,
	}, nil

}

// Uppercase ...
func (stringService) Uppercase(s string) (string, error) {

	if numb, err := strconv.Atoi(s); err == nil {
		return s, errs.New(err, fmt.Sprintf("do you want to do some Math with %d??", numb), http.StatusBadRequest)
	}
	return strings.ToUpper(s), nil
}

// Count ...
func (stringService) Count(s string) (int, error) {
	if s == "" {
		return 0, errs.New(errs.ErrEmpty, "why would want to count some  –e m p t y   s t r i n g– huh!??", http.StatusBadRequest)
	}
	return utf8.RuneCountInString(s), nil
}
