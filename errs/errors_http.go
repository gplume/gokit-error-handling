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
func EncodeError(logger kitlog.Logger, fullStack bool) kithttp.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		code := http.StatusInternalServerError
		httperr := []interface{}{
			"http.url", ctx.Value(kithttp.ContextKeyRequestURI),
			"http.method", ctx.Value(kithttp.ContextKeyRequestMethod),
			"http.user_agent", ctx.Value(kithttp.ContextKeyRequestUserAgent),
			"http.proto", ctx.Value(kithttp.ContextKeyRequestProto),
			// following is unneeded? why would we want only a part of the url anyway?
			// "http.path", ctx.Value(kithttp.ContextKeyRequestPath),
		}
		switch err {
		// case ErrInvalidBody:
		// 	code = http.StatusBadRequest
		// case sql.ErrNoRows:
		// 	code = http.StatusNotFound
		case err.(*Error):
			// fmt.Println("-----from errs.pkg-----")
			obj := []interface{}{
				"caller", err.(*Error).Caller,
				"message", err.(*Error).Message,
				"error", err.(*Error).Err,
				"code", err.(*Error).Code,
				"level", err.(*Error).Level,
			}
			if (fullStack) || err.(*Error).Caller == "" && err.Error() != "" {
				obj = append(obj, "stack", err.Error())
			}
			httperr = append(httperr, obj...)
			logger.Log(httperr...)
			if errCode := err.(*Error).Code; errCode > 0 {
				code = errCode
			}

		default:
			// fmt.Println("-----std.errors-----")
			obj := []interface{}{
				"error", err.Error(),
				"stack", fmt.Sprintf("%+v", err),
			}
			httperr = append(httperr, obj...)
			logger.Log(httperr...)
		}

		// Now we print to the client:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(code)
		// defaults the error message content
		var msg string
		// if from errs.pkg then retreive Message if not empty
		// or do some specific standardization message sorting
		if er, itis := err.(*Error); itis && er.Message != "" {
			msg = er.Message
		}
		// this should be removed for more granularity:
		// switch code {
		// case http.StatusBadRequest:
		// 	msg = ErrInvalidBody.Error()
		// case http.StatusInternalServerError:
		// 	msg = ErrInternalServer.Error()
		// 	// ...etc
		// }

		// response
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": fmt.Sprintf("%s", msg),
		})
	}
}
