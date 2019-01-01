package api

import (
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
		return s, errs.New(errs.ErrEmpty)
	}
	return vs.next.Uppercase(s)
}

func (vs serviceValidation) Count(s string) (int, error) {
	return vs.next.Count(s)
}
