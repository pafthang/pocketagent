package space

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	mwctx "github.com/pafthang/pocketagent/pkgs/middle/context"
)

func TestRequireHeaderOnly(t *testing.T) {
	e := echo.New()
	e.Use(Require(Options{}))
	e.GET("/", func(c echo.Context) error {
		id, ok := mwctx.SpaceIDFromContext(c)
		if !ok || id != "space-1" {
			t.Fatalf("unexpected space id: %q %v", id, ok)
		}
		return c.NoContent(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(mwctx.HeaderSpaceID, "space-1")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestRequireQueryFallback(t *testing.T) {
	e := echo.New()
	e.Use(Require(Options{AllowQueryFallback: true}))
	e.GET("/", func(c echo.Context) error {
		id, ok := mwctx.SpaceIDFromContext(c)
		if !ok || id != "space-2" {
			t.Fatalf("unexpected space id: %q %v", id, ok)
		}
		return c.NoContent(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/?space_id=space-2", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d", rec.Code)
	}
}