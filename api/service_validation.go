package api

import (
	"net/http"

	"github.com/gplume/gokit-error-handling/errs"
)

type serviceValidation struct {
	next Service
}

func newValidationService(svc Service) Service {
	return serviceValidation{
		next: svc,
	}
}

func (vs serviceValidation) Uppercase(s string) (string, error) {
	if s == "" {
		return s, errs.New(errs.ErrEmptyParam)
	}
	return vs.next.Uppercase(s)
}

func (vs serviceValidation) Count(s string) (int, error) {
	if s == "" {
		return 0, errs.New("why would want to count some  –e m p t y   s t r i n g– huh!??", http.StatusBadRequest)
	}
	return vs.next.Count(s)
}
