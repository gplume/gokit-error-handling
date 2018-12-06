package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
)

const (
	componentName = "errors-handling"
)

var (
	errInternalServer = errors.New("an internal server error occurred please contact the server's administrator")
	errInvalidBody    = errors.New("invalid body")
)

func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
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
	V   string `json:"v"`
	Err string `json:"err,omitempty"`
}

type countRequest struct {
	S string `json:"s"`
}

type countResponse struct {
	V int `json:"v"`
}

// ---------------------------------------------------------------------------------------------------------

// MakeHTTPHandler returns all http handler for the user service
func MakeHTTPHandler(
	endpoints Endpoints,
	logger kitlog.Logger,
) http.Handler {

	logger = kitlog.With(logger, "component", componentName)
	options := []kithttp.ServerOption{
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerErrorEncoder(encodeError(logger)),
	}

	uppercaseHandler := kithttp.NewServer(
		endpoints.Uppercase,
		decodeUppercaseRequest,
		encodeResponse,
		options...,
	)

	countHandler := kithttp.NewServer(
		endpoints.Count,
		decodeCountRequest,
		encodeResponse,
		options...,
	)

	r := mux.NewRouter().StrictSlash(true)
	{
		r.Handle("/uppercase", uppercaseHandler).Methods(http.MethodPost)
		r.Handle("/count", countHandler).Methods(http.MethodPost)
	}
	return RecoverFromPanic(logger, r)
}

func encodeError(logger kitlog.Logger) kithttp.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		code := http.StatusInternalServerError
		if sc, ok := err.(kithttp.StatusCoder); ok {
			fmt.Println("kithttp.StatusCoder:", sc)
			code = sc.StatusCode()
		}
		switch err {
		default:
			logger.Log(
				"err", fmt.Sprintf("%+v", err),
				"http.url", ctx.Value(kithttp.ContextKeyRequestURI),
				"http.path", ctx.Value(kithttp.ContextKeyRequestPath),
				"http.method", ctx.Value(kithttp.ContextKeyRequestMethod),
				"http.user_agent", ctx.Value(kithttp.ContextKeyRequestUserAgent),
				"http.proto", ctx.Value(kithttp.ContextKeyRequestProto),
			)
		case errInvalidBody:
			code = http.StatusBadRequest
		case sql.ErrNoRows:
			code = http.StatusNotFound
		}
		fmt.Println(err)

		w.WriteHeader(code)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
	}
}

// RecoverFromPanic is the Global recoverer in case of Panic.
func RecoverFromPanic(logger kitlog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Log("err", rec,
					"http.url", r.RequestURI,
					"http.path", r.URL.Path,
					"http.method", r.Method,
					"http.user_agent", r.Header.Get("User-Agent"),
					"http.proto", r.Proto)
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error": errInternalServer.Error(),
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
