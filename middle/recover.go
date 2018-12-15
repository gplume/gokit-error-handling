package middle

import (
	"encoding/json"
	"fmt"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
)

// RecoverFromPanic as main recover for the global Handler...
func RecoverFromPanic(logger kitlog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Log(
					"err_recoverFromPanic", rec,
					"http.url", r.RequestURI,
					"http.path", r.URL.Path,
					"http.method", r.Method,
					"http.user_agent", r.Header.Get("User-Agent"),
					"http.proto", r.Proto,
				)
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error": fmt.Sprintf("%v, %T", rec, rec),
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
