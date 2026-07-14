package taskapis

import (
	"github.com/labstack/echo/v4"
	natsclient "github.com/pafthang/pocketagent/internal/nats/client"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/models"
)

// RegisterRoutes wires task, schedule, and streaming endpoints.
func RegisterRoutes(
	tenant *echo.Group,
	nc *natsclient.Client,
	pb *pbclient.Client,
	taskRead, taskWrite echo.MiddlewareFunc,
) {
	toolCfg := defaultToolConfig()
	tenant.POST("/tasks", func(c echo.Context) error {
		var task models.Task
		if err := c.Bind(&task); err != nil {
			return err
		}
		return PublishTaskWithTools(c, nc, pb, task, toolCfg)
	}, taskWrite)
	tenant.GET("/tasks", listTasksHandler(pb), taskRead)
	tenant.GET("/tasks/:id/stream", func(c echo.Context) error {
		return sseTaskStream(c, nc, pb)
	}, taskRead)
	tenant.GET("/tasks/:id", getTaskHandler(pb), taskRead)
	tenant.DELETE("/tasks/:id", deleteTaskHandler(pb), taskWrite)
	tenant.GET("/ws/task/:taskId", func(c echo.Context) error {
		return wsTaskStream(c, nc, pb)
	}, taskRead)

	tenant.POST("/schedules", createScheduleHandler(pb), taskWrite)
	tenant.GET("/schedules", listSchedulesHandler(pb), taskRead)
	tenant.GET("/schedules/:id", getScheduleHandler(pb), taskRead)
	tenant.PATCH("/schedules/:id", updateScheduleHandler(pb), taskWrite)
	tenant.DELETE("/schedules/:id", deleteScheduleHandler(pb), taskWrite)
}
