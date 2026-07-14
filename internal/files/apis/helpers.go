package fileapis

import (
	"net/http"
	"strings"
	"time"

	filepath "github.com/pafthang/pocketagent/internal/files/path"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/models"
)

func resolveParentFolder(pb *pbclient.Client, spaceID, projectID, dirPath string) (string, error) {
	dirPath = filepath.NormalizePath(dirPath)
	if dirPath == "/" {
		return "", nil
	}
	if projectID != "" && dirPath == filepath.ProjectRoot(projectID) {
		return "", nil
	}
	folder, err := pb.FindFileByPath(spaceID, dirPath)
	if err != nil {
		return "", err
	}
	if !folder.IsDir {
		return "", &pbclient.APIError{StatusCode: http.StatusBadRequest, Message: "path is not a directory"}
	}
	if projectID != "" && folder.ProjectID != projectID {
		return "", &pbclient.APIError{StatusCode: http.StatusNotFound, Message: "folder not found"}
	}
	return folder.ID, nil
}

func validateProjectInSpace(pb *pbclient.Client, spaceID, projectID string) error {
	project, err := pb.GetProject(projectID)
	if err != nil {
		return err
	}
	if project.SpaceID != spaceID {
		return &pbclient.APIError{StatusCode: http.StatusNotFound, Message: "project not found"}
	}
	return nil
}

func loadFileInSpace(pb *pbclient.Client, spaceID, id string) (models.StoredFile, error) {
	file, err := pb.GetFile(id)
	if err != nil {
		return models.StoredFile{}, err
	}
	if file.SpaceID != spaceID {
		return models.StoredFile{}, &pbclient.APIError{StatusCode: http.StatusNotFound, Message: "file not found"}
	}
	return file, nil
}

func toBrowseEntry(file models.StoredFile) models.BrowseEntry {
	return models.BrowseEntry{
		ID:         file.ID,
		Name:       file.Name,
		Path:       file.VirtualPath,
		IsDir:      file.IsDir,
		Size:       file.Size,
		MimeType:   file.MimeType,
		ProjectID:  file.ProjectID,
		ModifiedAt: file.UpdatedAt,
	}
}

func memoDocumentIDForFile(fileID string) string {
	return "file-" + strings.TrimSpace(fileID)
}

func parsePBTime(raw string) int64 {
	if raw == "" {
		return time.Now().UnixMilli()
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02 15:04:05.000Z", "2006-01-02 15:04:05Z07:00"} {
		if ts, err := time.Parse(layout, raw); err == nil {
			return ts.UnixMilli()
		}
	}
	return time.Now().UnixMilli()
}