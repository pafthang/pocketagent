package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// RequestsTotal counts HTTP requests by service, method, and status.
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pocketagent_requests_total",
			Help: "Total number of requests",
		},
		[]string{"service", "method", "status"},
	)

	// TaskDuration records task execution latency.
	TaskDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "pocketagent_task_duration_seconds",
			Help: "Task execution duration",
		},
		[]string{"service", "status"},
	)

	// DLQMessagesTotal counts archived dead-letter messages.
	DLQMessagesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pocketagent_dlq_messages_total",
			Help: "Dead-letter queue messages archived after handler exhaustion",
		},
		[]string{"service", "reason"},
	)
)

func init() {
	prometheus.MustRegister(RequestsTotal)
	prometheus.MustRegister(TaskDuration)
	prometheus.MustRegister(DLQMessagesTotal)
}

// Handler returns the Prometheus scrape handler.
func Handler() http.Handler {
	return promhttp.Handler()
}