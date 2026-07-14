package fileapis

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes wires file storage HTTP endpoints.
func RegisterRoutes(tenant *echo.Group, deps *Deps, readAction, writeAction echo.MiddlewareFunc) {
	tenant.GET("/files/browse", browseFilesHandler(deps), readAction)
	tenant.GET("/files/recent", recentFilesHandler(deps), readAction)
	tenant.GET("/files/:id", getFileHandler(deps), readAction)
	tenant.GET("/files/:id/download", downloadFileHandler(deps), readAction)
	tenant.GET("/files/:id/content", fileContentHandler(deps), readAction)
	tenant.POST("/files/upload", uploadFileHandler(deps), writeAction)
	tenant.POST("/files/folders", createFolderHandler(deps), writeAction)
	tenant.POST("/files/:id/ingest", ingestFileHandler(deps), writeAction)
	tenant.DELETE("/files/:id", deleteFileHandler(deps), writeAction)
}