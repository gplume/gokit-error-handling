package api

import (
	"errors"
	"net/http"

	"github.com/gplume/gokit-error-handling/errs"
)

type serviceValidation struct {
	next Service
}

func newValidationService(svc Service, logger Logger) (Service, error) {
	return serviceValidation{
		next: svc,
	}, nil
}

func (vs serviceValidation) Uppercase(s string) (string, error) {
	if s == "" {
		return s, errs.New(errors.New("string is empty"), "...e m p t y   s t r i n g...", http.StatusBadRequest)
	}
	return vs.next.Uppercase(s)
}

func (vs serviceValidation) Count(s string) (int, error) {
	return vs.next.Count(s)
}
