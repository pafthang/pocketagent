package common

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics holds Prometheus metrics
var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pocketagent_requests_total",
			Help: "Total number of requests",
		},
		[]string{"service", "method", "status"},
	)

	TaskDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "pocketagent_task_duration_seconds",
			Help: "Task execution duration",
		},
		[]string{"service", "status"},
	)
)

func init() {
	prometheus.MustRegister(RequestsTotal)
	prometheus.MustRegister(TaskDuration)
}

// MetricsHandler returns Prometheus handler
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
