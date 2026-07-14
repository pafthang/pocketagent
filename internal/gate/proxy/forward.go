package proxy

import (
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	apimw "github.com/pafthang/pocketagent/pkgs/middle"
)

func bearerToken(c echo.Context, requireAuth bool) (string, bool) {
	if !requireAuth {
		return "", true
	}
	token := apimw.ExtractBearer(c)
	if token == "" {
		_ = c.JSON(http.StatusUnauthorized, map[string]string{"error": "authorization required"})
		return "", false
	}
	return token, true
}

func readBody(c echo.Context) ([]byte, error) {
	return io.ReadAll(c.Request().Body)
}

func spaceID(c echo.Context) string {
	id := strings.TrimSpace(c.Request().Header.Get(apimw.HeaderSpaceID))
	if id == "" {
		if fromCtx, ok := apimw.SpaceIDFromContext(c); ok {
			id = fromCtx
		}
	}
	return id
}

func targetPath(base, suffix, rawQuery string) string {
	target := base
	if suffix != "" {
		target = base + "/" + suffix
	}
	if rawQuery != "" {
		target += "?" + rawQuery
	}
	return target
}

func writeResponse(c echo.Context, resp *http.Response, extraHeaders ...map[string]string) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	for k, vals := range resp.Header {
		for _, v := range vals {
			c.Response().Header().Add(k, v)
		}
	}
	for _, hdrs := range extraHeaders {
		for k, v := range hdrs {
			c.Response().Header().Set(k, v)
		}
	}
	return c.Blob(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

func badGateway(c echo.Context, err error) error {
	return c.JSON(http.StatusBadGateway, map[string]string{"error": err.Error()})
}