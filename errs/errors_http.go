package errs

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
)

// EncodeError ...
func EncodeError(logger kitlog.Logger) kithttp.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		code := http.StatusInternalServerError
		switch err {
		case err.(*Error):
			fmt.Println("-----errs.errrorrrr-----")
			logger.Log(
				"caller", err.(*Error).Caller,
				"message", err.(*Error).Message,
				"error", err.(*Error).Err,
				"code", fmt.Sprintf("%v", err.(*Error).Code),
				"http.url", ctx.Value(kithttp.ContextKeyRequestURI),
				"http.path", ctx.Value(kithttp.ContextKeyRequestPath),
				"http.method", ctx.Value(kithttp.ContextKeyRequestMethod),
				"http.user_agent", ctx.Value(kithttp.ContextKeyRequestUserAgent),
				"http.proto", ctx.Value(kithttp.ContextKeyRequestProto),
				// "stack", err.Error(),
			)
			if errCode := err.(*Error).Code; errCode > 0 {
				code = errCode
			}
		default:
			fmt.Println("-----errrorrrr-----")
			logger.Log(
				"error", err.Error(),
				"http.url", ctx.Value(kithttp.ContextKeyRequestURI),
				"http.path", ctx.Value(kithttp.ContextKeyRequestPath),
				"http.method", ctx.Value(kithttp.ContextKeyRequestMethod),
				"http.user_agent", ctx.Value(kithttp.ContextKeyRequestUserAgent),
				"http.proto", ctx.Value(kithttp.ContextKeyRequestProto),
				"stack", fmt.Sprintf("%+v", err),
			)
		case ErrInvalidBody:
			code = http.StatusBadRequest
		case sql.ErrNoRows:
			code = http.StatusNotFound
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(code)
		msg := err.Error()
		if er, isit := err.(*Error); isit {
			msg = er.Message
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": fmt.Sprintf("%s", msg),
		})
	}
}
