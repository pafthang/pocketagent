module github.com/pafthang/pocketagent/services/agent-service

go 1.23

require (
	github.com/labstack/echo/v4 v4.13.3
	github.com/pafthang/pocketagent/internal/common v0.0.0
	github.com/pafthang/pocketagent/internal/models v0.0.0
	github.com/pafthang/pocketagent/internal/pocketbase v0.0.0
)

replace (
	github.com/pafthang/pocketagent/internal/common => ../../internal/common
	github.com/pafthang/pocketagent/internal/models => ../../internal/models
	github.com/pafthang/pocketagent/internal/pocketbase => ../../internal/pocketbase
)
