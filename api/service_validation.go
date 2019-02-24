package api

import (
	"net/http"

	"github.com/gplume/errs"
)

type serviceValidation struct {
	next Service
}

func newValidationService(svc Service) Service {
	return serviceValidation{
		next: svc,
	}
}

func (sv serviceValidation) Uppercase(s string) (string, error) {
	if s == "" {
		return s, errs.New(errs.ErrEmptyParam)
	}
	return sv.next.Uppercase(s)
}

func (sv serviceValidation) Count(s string) (int, error) {
	if s == "" {
		return 0, errs.New("why would want to count some  –e m p t y   s t r i n g– huh!??", http.StatusBadRequest)
	}
	return sv.next.Count(s)
}
