package proxy

import (
	"github.com/labstack/echo/v4"
	spaceclient "github.com/pafthang/pocketagent/internal/space/client"
)

// FixedPath forwards a single upstream path on the space service.
func FixedPath(sc *spaceclient.Client, method, path string, requireAuth bool) echo.HandlerFunc {
	return Space(sc, method, path, requireAuth)
}

// Space forwards wildcard paths under a base prefix to the space service.
func Space(sc *spaceclient.Client, method, path string, requireAuth bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := bearerToken(c, requireAuth)
		if !ok {
			return nil
		}

		body, err := readBody(c)
		if err != nil {
			return err
		}

		target := targetPath(path, c.Param("*"), "")

		resp, err := sc.Proxy(method, target, token, body)
		if err != nil {
			return badGateway(c, err)
		}
		return writeResponse(c, resp)
	}
}