package errs

import (
	"context"
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
		httperr := []interface{}{
			"http.url", ctx.Value(kithttp.ContextKeyRequestURI),
			"http.method", ctx.Value(kithttp.ContextKeyRequestMethod),
			"http.user_agent", ctx.Value(kithttp.ContextKeyRequestUserAgent),
			"http.proto", ctx.Value(kithttp.ContextKeyRequestProto),
		}
		switch err {
		case err.(*Err):
			// fmt.Println("-----from errs.pkg-----")
			if e, ok := err.(*Err); ok && err.(*Err).Level < startLoggingUnderLevel {
				obj := []interface{}{
					"caller", e.Caller,
					"message", e.Message,
					"error", e.Err,
					"code", e.Code,
					"level", e.Level.String(),
				}
				if e.Stack != nil {
					obj = append(obj, "stack", e.Error())
				}
				obj = append(obj, httperr...)
				logger.Log(obj...)
				if errCode := e.Code; errCode > 0 {
					code = errCode
				}
			}
		default:
			// fmt.Println("-----std.errors-----") // for backward compatibilty with std error
			obj := []interface{}{
				"error", err.Error(),
				"stack", fmt.Sprintf("%+v", err),
			}
			obj = append(obj, httperr...)
			logger.Log(obj...)
		}

		// Now we print to the client:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(code)
		// defaults the error message content
		var msg string
		// if from errs.pkg then retreive Message if not empty
		// but we should probably set an option for that here
		// for displaying std message if wished
		if er, itis := err.(*Err); itis && er.Message != "" {
			msg = er.Message
		}
		// in case of...
		if msg == "" {
			msg = ErrInternalServer.Message
		}

		// response to the client
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": fmt.Sprintf("%s", msg),
		})
	}
}
