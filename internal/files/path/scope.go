package path

import "strings"

// FileScope is a resolved browse/upload directory within a space.
type FileScope struct {
	ProjectID string
	DirPath   string
}

// ResolveScope maps an optional project_id and raw path into a virtual directory.
func ResolveScope(routeProjectID, rawPath, queryProjectID string) FileScope {
	projectID := strings.TrimSpace(queryProjectID)
	if pid := strings.TrimSpace(routeProjectID); pid != "" {
		projectID = pid
	}

	rawPath = strings.TrimSpace(rawPath)
	if pid, dir := ParseBrowsePath(rawPath); pid != "" {
		return FileScope{ProjectID: pid, DirPath: dir}
	}

	dirPath := NormalizePath(rawPath)
	if projectID == "" {
		return FileScope{DirPath: dirPath}
	}

	root := ProjectRoot(projectID)
	if dirPath == "/" {
		return FileScope{ProjectID: projectID, DirPath: root}
	}
	if dirPath == root || strings.HasPrefix(dirPath, root+"/") {
		return FileScope{ProjectID: projectID, DirPath: dirPath}
	}

	rel := strings.TrimPrefix(dirPath, "/")
	return FileScope{ProjectID: projectID, DirPath: JoinPath(root, rel)}
}

// BuildProjectPath joins a project virtual root with a relative segment.
func BuildProjectPath(projectID, relative string) string {
	return ResolveScope(projectID, relative, "").DirPath
}