package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/gplume/gokit-error-handling/errs"
	"github.com/gplume/gokit-error-handling/middle"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
)

const (
	componentName = "gokit-error-handling"
)

func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, errs.New("error decoding uppercase request", err, errs.Info)
	}
	return request, nil
}

func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, errs.New(err, "error decoding count request", http.StatusBadRequest, errs.Info)
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

type uppercaseRequest struct {
	S string `json:"string"`
}

type uppercaseResponse struct {
	V   string `json:"uppercased"`
	Err string `json:"err,omitempty"`
}

type countRequest struct {
	S string `json:"string"`
}

type countResponse struct {
	V int `json:"counted"`
}

// MakeHTTPHandler returns all http handler for the user service
func MakeHTTPHandler(
	endpoints Endpoints,
	logger kitlog.Logger,
) http.Handler {
	options := []kithttp.ServerOption{
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerErrorEncoder(errs.EncodeError(logger)),
	}

	uppercaseHandler := middle.Ware(kithttp.NewServer(
		endpoints.Uppercase,
		decodeUppercaseRequest,
		encodeResponse,
		options...,
	),
		middle.Notify(),
		middle.Metrics(componentName, "uppercase-handler"),
	)

	countHandler := middle.Ware(kithttp.NewServer(
		endpoints.Count,
		decodeCountRequest,
		encodeResponse,
		options...,
	),
		middle.Notify(),
		middle.Metrics(componentName, "count-handler"),
	)

	/*************** Bone Muxer *****************/
	router := bone.New()
	router.Post("/uppercase", uppercaseHandler)
	router.Post("/count", countHandler)

	return router
}
