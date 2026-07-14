package fileapis

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	filepath "github.com/pafthang/pocketagent/internal/files/path"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/httpx"
	apimw "github.com/pafthang/pocketagent/pkgs/middle"
	"github.com/pafthang/pocketagent/pkgs/models"
)

func uploadFileHandler(deps *Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}

		src, err := c.FormFile("file")
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "file is required"})
		}

		projectID := strings.TrimSpace(c.FormValue("project_id"))
		if projectID == "" {
			projectID = strings.TrimSpace(c.QueryParam("project_id"))
		}
		scope := filepath.ResolveScope("", c.FormValue("path"), projectID)
		if scope.ProjectID != "" {
			if err := validateProjectInSpace(deps.PB, spaceID, scope.ProjectID); err != nil {
				return httpx.MapPocketError(c, err)
			}
		}

		name := strings.TrimSpace(src.Filename)
		if err := filepath.ValidateName(name); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		virtualPath := filepath.JoinPath(scope.DirPath, name)
		if _, err := deps.PB.FindFileByPath(spaceID, virtualPath); err == nil {
			return c.JSON(http.StatusConflict, map[string]string{"error": "file already exists"})
		} else {
			var apiErr *pbclient.APIError
			if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusNotFound {
				return httpx.MapPocketError(c, err)
			}
		}

		parentID, err := resolveParentFolder(deps.PB, spaceID, scope.ProjectID, scope.DirPath)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		recordID := fmt.Sprintf("file-%d", time.Now().UnixNano())
		in, err := src.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		defer in.Close()

		storageKey, checksum, size, err := deps.Store.Save(spaceID, recordID, in)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		mimeType := filepath.DetectMimeType(name, nil)
		if header := src.Header.Get("Content-Type"); header != "" && header != "application/octet-stream" {
			mimeType = header
		}

		file := models.StoredFile{
			SpaceID:     spaceID,
			ProjectID:   scope.ProjectID,
			ParentID:    parentID,
			Name:        name,
			VirtualPath: virtualPath,
			MimeType:    mimeType,
			Size:        size,
			StorageKey:  storageKey,
			Checksum:    checksum,
		}
		if user, ok := apimw.UserFromContext(c); ok {
			file.UploadedBy = user.ID
		}

		stored, err := deps.PB.CreateFile(file)
		if err != nil {
			_ = deps.Store.Delete(storageKey)
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusCreated, stored)
	}
}

func createFolderHandler(deps *Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}

		var req createFolderRequest
		if err := c.Bind(&req); err != nil {
			return err
		}
		if err := filepath.ValidateName(req.Name); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		projectID := strings.TrimSpace(req.ProjectID)
		if projectID == "" {
			projectID = strings.TrimSpace(c.QueryParam("project_id"))
		}
		scope := filepath.ResolveScope("", req.Path, projectID)
		if scope.ProjectID != "" {
			if err := validateProjectInSpace(deps.PB, spaceID, scope.ProjectID); err != nil {
				return httpx.MapPocketError(c, err)
			}
		}

		virtualPath := filepath.JoinPath(scope.DirPath, req.Name)
		if _, err := deps.PB.FindFileByPath(spaceID, virtualPath); err == nil {
			return c.JSON(http.StatusConflict, map[string]string{"error": "path already exists"})
		} else {
			var apiErr *pbclient.APIError
			if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusNotFound {
				return httpx.MapPocketError(c, err)
			}
		}

		parentID, err := resolveParentFolder(deps.PB, spaceID, scope.ProjectID, scope.DirPath)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		folder := models.StoredFile{
			SpaceID:     spaceID,
			ProjectID:   scope.ProjectID,
			ParentID:    parentID,
			Name:        req.Name,
			VirtualPath: virtualPath,
			IsDir:       true,
		}
		if user, ok := apimw.UserFromContext(c); ok {
			folder.UploadedBy = user.ID
		}

		stored, err := deps.PB.CreateFile(folder)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusCreated, stored)
	}
}