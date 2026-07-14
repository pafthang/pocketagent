package common

import (
	"net/http"

	"github.com/pafthang/pocketagent/pkgs/common/metrics"
)

var (
	RequestsTotal    = metrics.RequestsTotal
	TaskDuration     = metrics.TaskDuration
	DLQMessagesTotal = metrics.DLQMessagesTotal
)

func MetricsHandler() http.Handler { return metrics.Handler() }