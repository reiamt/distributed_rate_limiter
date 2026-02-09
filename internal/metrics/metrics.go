package metrics

import "github.com/prometheus/client_golang/prometheus"

var RequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "ratelimiter_requests_total",
		Help: "Total number of rate limit checks",
	},
	[]string{"status"}, // "allowed" or "blocked"
)

var RequestDuration = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Name:    "ratelimiter_request_duration_seconds",
		Help:    "Duration of the rate limit checks",
		Buckets: prometheus.DefBuckets,
	},
)

var RedisErrorsTotal = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "ratelimiter_redis_errors_total",
		Help: "Total number of errors from redis",
	},
)

func init() {
	prometheus.MustRegister(RequestsTotal)
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(RedisErrorsTotal)
}
