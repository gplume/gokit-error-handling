package api

import (
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

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
	svc := Service(stringService{})
	svc = newValidationService(svc)
	return svc, nil
}

// Uppercase ...
func (stringService) Uppercase(s string) (string, error) {
	if _, err := strconv.Atoi(s); err == nil {
		return s, errs.New(http.StatusBadRequest, "uppercase some numbers dude, really??")
	}
	if s == "specific" {
		return s, errs.New(errs.ErrSpecific, errs.Low) // errs.Level overrides defined ErrSpecific.Level
	}
	if s == "specifics" {
		_, specErr := strconv.Atoi(s)
		return s, errs.New(specErr, errs.ErrInternalServer.Error())
	}
	return strings.ToUpper(s), nil
}

// Count ...
func (stringService) Count(s string) (int, error) {
	if s == "" {
		return 0, errs.New("why would want to count some  –e m p t y   s t r i n g– huh!??", http.StatusBadRequest)
	}
	return utf8.RuneCountInString(s), nil
}
