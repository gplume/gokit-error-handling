package middle

import (
	"fmt"
	"net/http"
	"time"

	kitlog "github.com/go-kit/kit/log"
)

// Notify ...
func Notify(logger kitlog.Logger) Wrapper {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Log("Notify( before )", fmt.Sprintf("%s", r.URL.Path))
			defer func(begin time.Time) {
				logger.Log("Notify( after )", time.Since(begin).Round(time.Nanosecond))
			}(time.Now())
			h.ServeHTTP(w, r)
		})
	}
}
