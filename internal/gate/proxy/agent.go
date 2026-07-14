package proxy

import (
	"github.com/labstack/echo/v4"
	agentclient "github.com/pafthang/pocketagent/internal/agent/client"
)

// Agent forwards wildcard paths under a base prefix to the agent service.
func Agent(ac *agentclient.Client, method, path string, requireAuth bool, responseHeaders ...map[string]string) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := bearerToken(c, requireAuth)
		if !ok {
			return nil
		}

		body, err := readBody(c)
		if err != nil {
			return err
		}

		target := targetPath(path, c.Param("*"), c.Request().URL.RawQuery)

		resp, err := ac.Proxy(method, target, token, spaceID(c), body)
		if err != nil {
			return badGateway(c, err)
		}
		return writeResponse(c, resp, responseHeaders...)
	}
}