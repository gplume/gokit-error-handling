package api

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints contains all go-kit like endpoints used to manipulate ratings
type Endpoints struct {
	Uppercase endpoint.Endpoint
	Count     endpoint.Endpoint
}

// MakeUppercaseEndpoint ...
func MakeUppercaseEndpoint(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(uppercaseRequest)
		v, err := svc.Uppercase(req.S)
		if err != nil {
			return nil, err
		}
		return uppercaseResponse{v, ""}, nil
	}
}

// MakeCountEndpoint ...
func MakeCountEndpoint(svc Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(countRequest)
		v, err := svc.Count(req.S)
		return countResponse{v}, err
	}
}
