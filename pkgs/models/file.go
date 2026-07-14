package models

// StoredFile is metadata for a space-scoped file or folder.
type StoredFile struct {
	ID           string   `json:"id"`
	SpaceID      string   `json:"space_id"`
	ProjectID    string   `json:"project_id,omitempty"`
	ParentID     string   `json:"parent_id,omitempty"`
	Name         string   `json:"name"`
	VirtualPath  string   `json:"virtual_path"`
	IsDir        bool     `json:"is_dir"`
	MimeType     string   `json:"mime_type,omitempty"`
	Size         int64    `json:"size,omitempty"`
	StorageKey   string   `json:"storage_key,omitempty"`
	Checksum     string   `json:"checksum,omitempty"`
	MemoIngested bool     `json:"memo_ingested,omitempty"`
	UploadedBy   string   `json:"uploaded_by,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	CreatedAt    string   `json:"created_at,omitempty"`
	UpdatedAt    string   `json:"updated_at,omitempty"`
}

// BrowseEntry is a single row returned by the file explorer.
type BrowseEntry struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	IsDir       bool   `json:"is_dir"`
	Size        int64  `json:"size,omitempty"`
	MimeType    string `json:"mime_type,omitempty"`
	ProjectID   string `json:"project_id,omitempty"`
	ModifiedAt  string `json:"modified_at,omitempty"`
}

// RecentFileEntry is a compact row for recent-file lists in the explorer.
type RecentFileEntry struct {
	Path      string `json:"path"`
	Name      string `json:"name"`
	IsDir     bool   `json:"is_dir"`
	Extension string `json:"extension"`
	Timestamp int64  `json:"timestamp"`
	Tool      string `json:"tool"`
}