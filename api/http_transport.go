package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/gplume/gokit-error-handling/handle"
	"github.com/gplume/gokit-error-handling/middle"
	"github.com/gplume/gokit-error-handling/utils"
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
		kithttp.ServerErrorEncoder(encodeError(logger)),
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

	/*************** PAT Muxer *****************/
	router := pat.New()
	{
		router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			utils.JSON(w, http.StatusOK, utils.JSMAP{"msg": fmt.Sprintf("route (%s) not found, sorry", r.URL.Path)})
		})
		router.Get("/", homeHandler)
		router.Get("/uppercase", http.HandlerFunc(theHome))
		router.Get("/:ppat", homeHandler)
		router.Post("/uppercase", uppercaseHandler)
		router.Post("/count", countHandler)
	}
	return router

	/*************** GORILLA/MUX **************/
	// router := mux.NewRouter().StrictSlash(true)
	// {
	// 	router.Handle("/uppercase", uppercaseHandler).Methods(http.MethodPost)
	// 	router.Handle("/count", countHandler).Methods(http.MethodPost)
	// }
	// return router

	/*************** NO-MUX *******************/
	// router := &handle.Handlers{
	// 	// route: "/"
	// 	HomeHandler: middle.Ware(new(handle.Home),
	// 		middle.Notify(logger),
	// 	),
	// 	// route: "/uppercase"
	// 	UpperCaseHandler: middle.Ware(
	// 		&handle.UpperCaseHandler{KitHandler: uppercaseHandler},
	// 		// middle.Notify(logger),
	// 	),
	// 	// route: "/count"
	// 	CharCountHandler: middle.Ware(
	// 		&handle.CharCountHandler{KitHandler: countHandler},
	// 		// middle.Notify(logger),
	// 	),
	// }
	// return router
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
			"error": fmt.Sprintf("%+v", err),
		})
	}
}

// RecoverFromPanic is the Global recoverer in case of Panic.
func RecoverFromPanic(logger kitlog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Log(
					"err", rec,
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

func theHome(w http.ResponseWriter, r *http.Request) {
	utils.JSON(w, http.StatusOK, utils.JSMAP{"msg": "HOME!"})
	return
}
