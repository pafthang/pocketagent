package fileapis

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	filepath "github.com/pafthang/pocketagent/internal/files/path"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/httpx"
	"github.com/pafthang/pocketagent/pkgs/models"
)

func patchFileHandler(deps *Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}
		file, err := loadFileInSpace(deps.PB, spaceID, c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		var req patchFileRequest
		if err := c.Bind(&req); err != nil {
			return err
		}

		if req.Name != "" {
			file, err = relocateFileRecord(deps.PB, file, req.Name, filepath.ParentPath(file.VirtualPath))
			if err != nil {
				return httpx.MapPocketError(c, err)
			}
		}
		if req.Tags != nil {
			file.Tags = req.Tags
			file, err = deps.PB.UpdateFile(file.ID, file)
			if err != nil {
				return httpx.MapPocketError(c, err)
			}
		}
		return c.JSON(http.StatusOK, file)
	}
}

func moveFileHandler(deps *Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		spaceID, ok := httpx.RequireSpaceID(c)
		if !ok {
			return nil
		}
		file, err := loadFileInSpace(deps.PB, spaceID, c.Param("id"))
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		var req pathTargetRequest
		if err := c.Bind(&req); err != nil {
			return err
		}
		destDir := filepath.NormalizePath(req.Path)
		file, err = relocateFileRecord(deps.PB, file, file.Name, destDir)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusOK, file)
	}
}

func copyFileHandler(deps *Deps) echo.HandlerFunc {
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
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "directory copy is not supported"})
		}
		if file.StorageKey == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "not a copyable file"})
		}

		var req pathTargetRequest
		if err := c.Bind(&req); err != nil {
			return err
		}
		destDir := filepath.NormalizePath(req.Path)
		if file.ProjectID != "" {
			if err := validateProjectInSpace(deps.PB, spaceID, file.ProjectID); err != nil {
				return httpx.MapPocketError(c, err)
			}
		}

		destPath := filepath.JoinPath(destDir, file.Name)
		if err := ensurePathAvailable(deps.PB, spaceID, destPath); err != nil {
			return httpx.MapPocketError(c, err)
		}

		parentID, err := resolveParentFolder(deps.PB, spaceID, file.ProjectID, destDir)
		if err != nil {
			return httpx.MapPocketError(c, err)
		}

		src, err := deps.Store.Open(file.StorageKey)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		defer src.Close()

		storageKey, checksum, size, err := deps.Store.Save(spaceID, fmt.Sprintf("%s-copy-%d", file.ID, time.Now().UnixNano()), src)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		copyRecord := models.StoredFile{
			SpaceID:     spaceID,
			ProjectID:   file.ProjectID,
			ParentID:    parentID,
			Name:        file.Name,
			VirtualPath: destPath,
			MimeType:    file.MimeType,
			Size:        size,
			StorageKey:  storageKey,
			Checksum:    checksum,
		}
		stored, err := deps.PB.CreateFile(copyRecord)
		if err != nil {
			_ = deps.Store.Delete(storageKey)
			return httpx.MapPocketError(c, err)
		}
		return c.JSON(http.StatusCreated, stored)
	}
}

func putFileContentHandler(deps *Deps) echo.HandlerFunc {
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
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "cannot write content to a directory"})
		}

		reader, mimeType, err := contentReaderFromRequest(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		oldKey := file.StorageKey
		storageKey, checksum, size, err := deps.Store.Save(spaceID, file.ID, reader)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		file.StorageKey = storageKey
		file.Checksum = checksum
		file.Size = size
		if mimeType != "" {
			file.MimeType = mimeType
		} else if file.MimeType == "" {
			file.MimeType = filepath.DetectMimeType(file.Name, nil)
		}

		updated, err := deps.PB.UpdateFile(file.ID, file)
		if err != nil {
			_ = deps.Store.Delete(storageKey)
			return httpx.MapPocketError(c, err)
		}
		if oldKey != "" && oldKey != storageKey {
			_ = deps.Store.Delete(oldKey)
		}
		return c.JSON(http.StatusOK, updated)
	}
}

func contentReaderFromRequest(c echo.Context) (io.Reader, string, error) {
	ct := strings.ToLower(strings.TrimSpace(c.Request().Header.Get("Content-Type")))
	if strings.HasPrefix(ct, "application/json") {
		var body putFileContentJSON
		if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
			return nil, "", err
		}
		return strings.NewReader(body.Content), "text/plain", nil
	}
	mime := ct
	if mime == "" {
		mime = "application/octet-stream"
	}
	return c.Request().Body, mime, nil
}

func relocateFileRecord(pb *pbclient.Client, file models.StoredFile, name, destDir string) (models.StoredFile, error) {
	name = strings.TrimSpace(name)
	if err := filepath.ValidateName(name); err != nil {
		return models.StoredFile{}, err
	}
	destDir = filepath.NormalizePath(destDir)
	newPath := filepath.JoinPath(destDir, name)
	if newPath == file.VirtualPath {
		return file, nil
	}
	if err := ensurePathAvailable(pb, file.SpaceID, newPath); err != nil {
		return models.StoredFile{}, err
	}

	parentID, err := resolveParentFolder(pb, file.SpaceID, file.ProjectID, destDir)
	if err != nil {
		return models.StoredFile{}, err
	}

	oldPath := file.VirtualPath
	file.Name = name
	file.VirtualPath = newPath
	file.ParentID = parentID

	if file.IsDir {
		descendants, err := collectDescendants(pb, file.SpaceID, file.ID, file.ProjectID)
		if err != nil {
			return models.StoredFile{}, err
		}
		if _, err := pb.UpdateFile(file.ID, file); err != nil {
			return models.StoredFile{}, err
		}
		for _, child := range descendants {
			child.VirtualPath = rebaseVirtualPath(oldPath, newPath, child.VirtualPath)
			if _, err := pb.UpdateFile(child.ID, child); err != nil {
				return models.StoredFile{}, err
			}
		}
		return pb.GetFile(file.ID)
	}

	return pb.UpdateFile(file.ID, file)
}

func ensurePathAvailable(pb *pbclient.Client, spaceID, virtualPath string) error {
	_, err := pb.FindFileByPath(spaceID, virtualPath)
	if err == nil {
		return &pbclient.APIError{StatusCode: http.StatusConflict, Message: "path already exists"}
	}
	var apiErr *pbclient.APIError
	if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
		return nil
	}
	return err
}

func collectDescendants(pb *pbclient.Client, spaceID, folderID, projectID string) ([]models.StoredFile, error) {
	children, _, err := pb.ListChildren(spaceID, folderID, projectID, 1, 500)
	if err != nil {
		return nil, err
	}
	out := make([]models.StoredFile, 0, len(children))
	for _, child := range children {
		out = append(out, child)
		if child.IsDir {
			sub, err := collectDescendants(pb, spaceID, child.ID, projectID)
			if err != nil {
				return nil, err
			}
			out = append(out, sub...)
		}
	}
	return out, nil
}

func rebaseVirtualPath(oldRoot, newRoot, path string) string {
	if path == oldRoot {
		return newRoot
	}
	prefix := oldRoot + "/"
	if strings.HasPrefix(path, prefix) {
		return newRoot + "/" + strings.TrimPrefix(path, prefix)
	}
	return path
}