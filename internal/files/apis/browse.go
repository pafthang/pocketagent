package fileapis

import (
	"net/http"

	"github.com/labstack/echo/v4"
	filepath "github.com/pafthang/pocketagent/internal/files/path"
	"github.com/pafthang/pocketagent/pkgs/httpx"
	"github.com/pafthang/pocketagent/pkgs/models"
)

func browseFilesHandler(deps *Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}

		scope := filepath.ResolveScope("", c.QueryParam("path"), c.QueryParam("project_id"))
		if scope.ProjectID != "" {
			if err := validateProjectInSpace(deps.PB, spaceID, scope.ProjectID); err != nil {
				return httpx.MapPocketError(c, err)
			}
		}

		parentID, err := resolveParentFolder(deps.PB, spaceID, scope.ProjectID, scope.DirPath)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		children, _, err := deps.PB.ListChildren(spaceID, parentID, scope.ProjectID, 1, 500)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		entries := make([]models.BrowseEntry, 0, len(children))
		for _, child := range children {
			entries = append(entries, toBrowseEntry(child))
		}

		return c.JSON(http.StatusOK, browseResponse{
			Path:      scope.DirPath,
			ProjectID: scope.ProjectID,
			Files:     entries,
		})
	}
}