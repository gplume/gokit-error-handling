package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/gplume/gokit-error-handling/errs"
	"github.com/gplume/gokit-error-handling/handle"
	"github.com/gplume/gokit-error-handling/middle"
	"github.com/gplume/gokit-error-handling/utils"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
)

const (
	componentName = "errors-handling"
)

func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, errs.New(err, "error decoding uppercase request", http.StatusBadRequest)

	}
	return request, nil
}

func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, errs.New(err, "error decoding count request", http.StatusBadRequest)
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

type uppercaseRequest struct {
	S string `json:"s"`
}

type uppercaseResponse struct {
	V   string `json:"uppercased"`
	Err string `json:"err,omitempty"`
}

type countRequest struct {
	S string `json:"s"`
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

	homeHandler := middle.Ware(new(handle.Home),
		middle.Notify(logger),
	)

	uppercaseHandler := middle.Ware(kithttp.NewServer(
		endpoints.Uppercase,
		decodeUppercaseRequest,
		encodeResponse,
		options...,
	),
		middle.Notify(logger),
	)

	countHandler := middle.Ware(kithttp.NewServer(
		endpoints.Count,
		decodeCountRequest,
		encodeResponse,
		options...,
	),
		middle.Notify(logger),
	)

	/*************** Bone Muxer *****************/
	router := bone.New()
	router.Get("/", homeHandler)
	router.GetFunc("/uppercase", uppercaseHome)
	router.Get("/:ppat", homeHandler)
	router.Post("/uppercase", uppercaseHandler)
	router.Post("/count", countHandler)
	return router
}

func uppercaseHome(w http.ResponseWriter, r *http.Request) {
	utils.JSON(w, http.StatusOK, utils.JSMAP{"msg": "UPPERCASE HOME!"})
	return
}
