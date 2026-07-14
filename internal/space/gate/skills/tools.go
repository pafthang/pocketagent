package skillapis

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pafthang/pocketagent/pkgs/httpx"
	"github.com/pafthang/pocketagent/internal/exec/tools"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
)

// ListToolsHandler returns tools available in the current space.
func ListToolsHandler(pb *pbclient.Client, toolCfg tools.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}

		result, err := tools.CollectSpaceTools(pb, spaceID, toolCfg)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusOK, map[string]interface{}{"tools": result})
	}
}
