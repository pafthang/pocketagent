package fileapis

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	filepath "github.com/pafthang/pocketagent/internal/files/path"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/httpx"
	"github.com/pafthang/pocketagent/pkgs/models"
)

func recentFilesHandler(deps *Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}

		limit := 20
		if n, err := strconv.Atoi(c.QueryParam("limit")); err == nil && n > 0 && n <= 100 {
			limit = n
		}

		projectID := strings.TrimSpace(c.QueryParam("project_id"))
		if projectID != "" {
			if err := validateProjectInSpace(deps.PB, spaceID, projectID); err != nil {
				return httpx.MapPocketError(c, err)
			}
		}

		filter := fmt.Sprintf("%s && is_dir = false", pbclient.FilesFilter(spaceID))
		if projectID != "" {
			filter += fmt.Sprintf(" && project_id = %q", projectID)
		} else {
			filter += ` && project_id = ""`
		}
		records, _, err := deps.PB.ListFiles(pbclient.ListOptions{Page: 1, PerPage: limit, Filter: filter})
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		recent := make([]models.RecentFileEntry, 0, len(records))
		for _, record := range records {
			recent = append(recent, models.RecentFileEntry{
				Path:      record.VirtualPath,
				Name:      record.Name,
				IsDir:     false,
				Extension: filepath.FileExtension(record.Name),
				Timestamp: parsePBTime(record.UpdatedAt),
				Tool:      "upload",
			})
		}

		return c.JSON(http.StatusOK, recentResponse{
			Files:     recent,
			ProjectID: projectID,
		})
	}
}