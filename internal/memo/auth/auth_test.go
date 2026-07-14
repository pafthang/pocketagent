package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestRequireServiceToken(t *testing.T) {
	e := echo.New()
	e.GET("/documents", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	}, RequireServiceToken("secret-token"))

	t.Run("valid header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/documents", nil)
		req.Header.Set(HeaderServiceToken, "secret-token")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/documents", nil)
		req.Header.Set(HeaderServiceToken, "wrong")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401", rec.Code)
		}
	})

	t.Run("empty expected skips auth", func(t *testing.T) {
		open := echo.New()
		open.GET("/documents", func(c echo.Context) error {
			return c.String(http.StatusOK, "ok")
		}, RequireServiceToken(""))

		req := httptest.NewRequest(http.MethodGet, "/documents", nil)
		rec := httptest.NewRecorder()
		open.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
	})
}