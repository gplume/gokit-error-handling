package middle

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpSuccessRegex, _ = regexp.Compile("^2[0-9]{2}$")
	// RequestCount ...
	RequestCount *kitprometheus.Counter
	// RequestLatency ...
	RequestLatency *kitprometheus.Histogram
	// CountResult ...
	CountResult *kitprometheus.Summary
)

func init() {

	// fieldKeys := []string{"method", "error"}
	RequestCount = kitprometheus.NewCounterFrom(prometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, []string{"component", "handler", "code", "method", "success"})
	RequestLatency = kitprometheus.NewHistogramFrom(prometheus.HistogramOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, []string{"component", "handler", "success"})
	// CountResult = kitprometheus.NewSummaryFrom(prometheus.SummaryOpts{
	// 	Namespace: "my_group",
	// 	Subsystem: "string_service",
	// 	Name:      "count_result",
	// 	Help:      "The result of each count method.",
	// }, []string{}) // no fields here
}

// Uppercase ...
// func (mw Instrumenting) Uppercase(s string) (output string, err error) {
// 	defer func(begin time.Time) {
// 		lvs := []string{"method", "uppercase", "error", fmt.Sprint(err != nil)}
// 		mw.RequestCount.With(lvs...).Add(1)
// 		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
// 	}(time.Now())

// 	output, err = mw.Next.Uppercase(s)
// 	return
// }

// // Count ...
// func (mw Instrumenting) Count(s string) (n int) {
// 	defer func(begin time.Time) {
// 		lvs := []string{"method", "count", "error", "false"}
// 		mw.RequestCount.With(lvs...).Add(1)
// 		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
// 		mw.CountResult.Observe(float64(n))
// 	}(time.Now())

// 	n = mw.Next.Count(s)
// 	return
// }

// ResponseWriter wraps the http.ResponseWriter for adding
// the https response status code
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader writes the http response status code and capture it inside the writer
func (lrw *ResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// NewResponseWriter implements the ResponseWriter interface and is used
// for capturing the http response status code
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, http.StatusOK}
}

// Metrics wraps a http handler for counting requests call and measuring request latency
func Metrics(componentName string, handlerName string) Wrapper {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lrw := NewResponseWriter(w)
			defer func(begin time.Time) {
				success := httpSuccessRegex.MatchString(strconv.Itoa(lrw.statusCode))
				RequestCount.With("component", componentName, "handler", handlerName, "code", strconv.Itoa(lrw.statusCode), "method", strings.ToLower(r.Method), "success", strconv.FormatBool(success)).Add(1)
				RequestLatency.With("component", componentName, "handler", handlerName, "success", strconv.FormatBool(success)).Observe(time.Since(begin).Seconds())
			}(time.Now())

			next.ServeHTTP(lrw, r)
		})
	}
}
