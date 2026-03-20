package stunserver

import "github.com/prometheus/client_golang/prometheus"

var (
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "stun_requests_total",
			Help: "Total STUN requests",
		},
		[]string{"type"},
	)

	ErrorCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "stun_errors_total",
			Help: "Total STUN errors",
		},
	)

	 RequestDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "stun_request_duration_seconds",
			Help: "Request latency",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func InitMetrics() {
	prometheus.MustRegister(RequestCount, ErrorCount, RequestDuration)
}
