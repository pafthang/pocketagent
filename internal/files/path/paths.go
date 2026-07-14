package path

import (
	"fmt"
	"path"
	"strings"
)

const projectPrefix = "/projects/"

// NormalizePath cleans a virtual browse path.
func NormalizePath(p string) string {
	p = strings.TrimSpace(p)
	if p == "" || p == "." {
		return "/"
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	p = path.Clean(p)
	if p == "." {
		return "/"
	}
	return p
}

// ParseBrowsePath extracts project scope and directory path.
func ParseBrowsePath(raw string) (projectID string, dirPath string) {
	dirPath = NormalizePath(raw)
	if !strings.HasPrefix(dirPath, projectPrefix) {
		return "", dirPath
	}
	rest := strings.TrimPrefix(dirPath, projectPrefix)
	parts := strings.SplitN(rest, "/", 2)
	if parts[0] == "" {
		return "", dirPath
	}
	projectID = parts[0]
	if len(parts) == 1 {
		return projectID, projectPrefix + projectID
	}
	return projectID, projectPrefix + projectID + "/" + parts[1]
}

// JoinPath joins a directory and entry name.
func JoinPath(dirPath, name string) string {
	dirPath = NormalizePath(dirPath)
	name = strings.TrimSpace(name)
	if name == "" {
		return dirPath
	}
	if dirPath == "/" {
		return "/" + name
	}
	return path.Clean(dirPath + "/" + name)
}

// ParentPath returns the parent directory of a virtual path.
func ParentPath(p string) string {
	p = NormalizePath(p)
	if p == "/" {
		return "/"
	}
	parent := path.Dir(p)
	if parent == "." {
		return "/"
	}
	return parent
}

// ProjectRoot returns the virtual root for a project.
func ProjectRoot(projectID string) string {
	return projectPrefix + strings.TrimSpace(projectID)
}

// ValidateName rejects invalid file or folder names.
func ValidateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("name is required")
	}
	if name == "." || name == ".." {
		return fmt.Errorf("invalid name")
	}
	if strings.ContainsAny(name, "/\\") {
		return fmt.Errorf("name cannot contain slashes")
	}
	return nil
}

// FileExtension returns lowercase extension without dot.
func FileExtension(name string) string {
	ext := strings.ToLower(strings.TrimPrefix(path.Ext(name), "."))
	return ext
}