package fileapis

import (
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	filepath "github.com/pafthang/pocketagent/internal/files/path"
	"github.com/pafthang/pocketagent/pkgs/httpx"
)

func getFileHandler(deps *Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}
		file, err := loadFileInSpace(deps.PB, spaceID, c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusOK, file)
	}
}

func downloadFileHandler(deps *Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}
		file, err := loadFileInSpace(deps.PB, spaceID, c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if file.IsDir || file.StorageKey == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "not a downloadable file"})
		}

		f, err := deps.Store.Open(file.StorageKey)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		defer f.Close()

		contentType := file.MimeType
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf(`attachment; filename="%s"`, file.Name))
		return c.Stream(http.StatusOK, contentType, f)
	}
}

func fileContentHandler(deps *Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}
		file, err := loadFileInSpace(deps.PB, spaceID, c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		if file.IsDir || file.StorageKey == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "not a file"})
		}
		if !filepath.IsTextMime(file.MimeType) {
			return c.JSON(http.StatusUnsupportedMediaType, map[string]string{"error": "binary file preview not supported"})
		}

		f, err := deps.Store.Open(file.StorageKey)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		defer f.Close()

		body, err := io.ReadAll(io.LimitReader(f, 2<<20))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, fileContentResponse{
			ID:      file.ID,
			Path:    file.VirtualPath,
			Content: string(body),
		})
	}
}

func deleteFileHandler(deps *Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}
		file, err := loadFileInSpace(deps.PB, spaceID, c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		if file.IsDir {
			children, _, err := deps.PB.ListChildren(spaceID, file.ID, file.ProjectID, 1, 1)
			if err != nil {
				return httpx.MapPocketError(c, err)
			}
			if len(children) > 0 {
				return c.JSON(http.StatusConflict, map[string]string{"error": "folder is not empty"})
			}
		} else if file.StorageKey != "" {
			_ = deps.Store.Delete(file.StorageKey)
		}

		if err := deps.PB.DeleteFile(file.ID); err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.NoContent(http.StatusNoContent)
	}
}