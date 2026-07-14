package fileapis

import "github.com/pafthang/pocketagent/pkgs/models"

type ingestFileRequest struct {
	Force bool     `json:"force"`
	Tags  []string `json:"tags"`
}

type createFolderRequest struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	ProjectID string `json:"project_id"`
}

type browseResponse struct {
	Path      string               `json:"path"`
	ProjectID string               `json:"project_id,omitempty"`
	Files     []models.BrowseEntry `json:"files"`
}

type recentResponse struct {
	Files     []models.RecentFileEntry `json:"files"`
	ProjectID string                   `json:"project_id,omitempty"`
}

type fileContentResponse struct {
	ID      string `json:"id"`
	Path    string `json:"path"`
	Content string `json:"content"`
}