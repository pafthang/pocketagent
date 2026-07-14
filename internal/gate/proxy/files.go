package proxy

import (
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	filesclient "github.com/pafthang/pocketagent/internal/files/client"
)

// Files forwards wildcard paths under /files to the files service.
func Files(fc *filesclient.Client, method, path string, requireAuth bool) echo.HandlerFunc {
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
		return forwardFiles(c, fc, method, target, token, body)
	}
}

// ProjectFiles rewrites /projects/:id/files/* to files service browse/upload routes.
func ProjectFiles(fc *filesclient.Client, method string, requireAuth bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := bearerToken(c, requireAuth)
		if !ok {
			return nil
		}

		body, err := readBody(c)
		if err != nil {
			return err
		}

		projectID := strings.TrimSpace(c.Param("id"))
		suffix := strings.Trim(strings.TrimSpace(c.Param("*")), "/")
		target := rewriteProjectFilesTarget(projectID, suffix, c.Request().URL.RawQuery)

		return forwardFiles(c, fc, method, target, token, body)
	}
}

func rewriteProjectFilesTarget(projectID, suffix, rawQuery string) string {
	var path string
	switch suffix {
	case "", "browse":
		path = "/files/browse"
	case "recent":
		path = "/files/recent"
	case "upload":
		path = "/files/upload"
	case "folders":
		path = "/files/folders"
	default:
		path = "/files/" + suffix
	}

	q := url.Values{}
	if rawQuery != "" {
		if parsed, err := url.ParseQuery(rawQuery); err == nil {
			for key, vals := range parsed {
				for _, val := range vals {
					q.Add(key, val)
				}
			}
		}
	}
	if projectID != "" {
		q.Set("project_id", projectID)
	}
	if encoded := q.Encode(); encoded != "" {
		return path + "?" + encoded
	}
	return path
}

func forwardFiles(c echo.Context, fc *filesclient.Client, method, target, token string, body []byte) error {
	resp, err := fc.Proxy(method, target, token, spaceID(c), body, c.Request().Header.Get("Content-Type"))
	if err != nil {
		return badGateway(c, err)
	}
	return writeResponse(c, resp)
}