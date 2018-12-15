package middle

import (
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// RequestCount ...
	RequestCount *kitprometheus.Counter
	// RequestLatency ...
	RequestLatency *kitprometheus.Histogram
	// CountResult ...
	CountResult *kitprometheus.Summary
)

func init() {

	fieldKeys := []string{"method", "error"}
	RequestCount = kitprometheus.NewCounterFrom(prometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	RequestLatency = kitprometheus.NewHistogramFrom(prometheus.HistogramOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	CountResult = kitprometheus.NewSummaryFrom(prometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{}) // no fields here
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
