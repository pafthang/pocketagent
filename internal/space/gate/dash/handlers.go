package dashboardapis

import (
	"github.com/pafthang/pocketagent/pkgs/httpx"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
)

func RegisterRoutes(tenant *echo.Group, pb *pbclient.Client, readAction echo.MiddlewareFunc) {
	tenant.GET("/dashboard", getDashboardHandler(pb), readAction)

	tenant.GET("/kits", listKitsHandler(pb), readAction)
	tenant.GET("/kits/catalog", listKitCatalogHandler(), readAction)
	tenant.GET("/kits/:id/data", getKitDataHandler(pb), readAction)
	tenant.POST("/kits/:id/activate", activateKitHandler(), readAction)
}

func getDashboardHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}
		limit := parseDashboardLimit(c.QueryParam("limit"), 25)

		summary, err := BuildSummary(pb, spaceID, BuildOptions{
			RecentLimit:   limit,
			ActivityLimit: limit,
		})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusOK, summary)
	}
}

func listKitsHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		if _, ok := httpx.RequireSpaceID(c); !ok {
			return nil
		}
		_ = pb
		return c.JSON(http.StatusOK, map[string]interface{}{
			"kits": []map[string]interface{}{BuiltinKit()},
		})
	}
}

func listKitCatalogHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		if _, ok := httpx.RequireSpaceID(c); !ok {
			return nil
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"catalog": []map[string]interface{}{BuiltinCatalogEntry()},
		})
	}
}

func getKitDataHandler(pb *pbclient.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}
		if c.Param("id") != BuiltinKitID {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "kit not found"})
		}

		limit := parseDashboardLimit(c.QueryParam("limit"), 50)
		summary, err := BuildSummary(pb, spaceID, BuildOptions{
			RecentLimit:   limit,
			ActivityLimit: limit,
		})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusOK, KitData(summary))
	}
}

func activateKitHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		if _, ok := httpx.RequireSpaceID(c); !ok {
			return nil
		}
		if c.Param("id") != BuiltinKitID {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "kit not found"})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"ok":     true,
			"id":     BuiltinKitID,
			"active": true,
		})
	}
}

func parseDashboardLimit(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n <= 0 {
		return fallback
	}
	if n > 200 {
		return 200
	}
	return n
}
