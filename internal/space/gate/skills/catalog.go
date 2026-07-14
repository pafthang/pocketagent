package skillapis

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/pafthang/pocketagent/pkgs/common"
)

type skillCatalogEntry struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Category     string   `json:"category"`
	ArgumentHint string   `json:"argument_hint"`
	Prompt       string   `json:"prompt"`
	Tools        []string `json:"tools"`
}

type skillCatalogResponse struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Category     string   `json:"category"`
	ArgumentHint string   `json:"argument_hint"`
	Prompt       string   `json:"prompt,omitempty"`
	Tools        []string `json:"tools,omitempty"`
	Installed    bool     `json:"installed"`
}

func loadSkillsCatalog(installed map[string]struct{}) ([]skillCatalogEntry, error) {
	dir, err := common.FindConfigsDir()
	if err != nil {
		return nil, err
	}

	raw, err := os.ReadFile(filepath.Join(dir, "skills-catalog.json"))
	if err != nil {
		return nil, err
	}

	var entries []skillCatalogEntry
	if err := json.Unmarshal(raw, &entries); err != nil {
		return nil, err
	}
	_ = installed
	return entries, nil
}

func catalogToResponses(entries []skillCatalogEntry, installed map[string]struct{}) []skillCatalogResponse {
	out := make([]skillCatalogResponse, 0, len(entries))
	for _, entry := range entries {
		_, isInstalled := installed[strings.ToLower(strings.TrimSpace(entry.Name))]
		out = append(out, skillCatalogResponse{
			ID:           entry.ID,
			Name:         entry.Name,
			Description:  entry.Description,
			Category:     entry.Category,
			ArgumentHint: entry.ArgumentHint,
			Prompt:       entry.Prompt,
			Tools:        entry.Tools,
			Installed:    isInstalled,
		})
	}
	return out
}

func searchCatalog(query string) ([]skillCatalogResponse, error) {
	query = strings.ToLower(strings.TrimSpace(query))
	entries, err := loadSkillsCatalog(nil)
	if err != nil {
		return nil, err
	}
	if query == "" {
		return catalogToResponses(entries, nil), nil
	}

	matches := make([]skillCatalogEntry, 0)
	for _, entry := range entries {
		haystack := strings.ToLower(entry.Name + " " + entry.Description + " " + entry.Category)
		if strings.Contains(haystack, query) {
			matches = append(matches, entry)
		}
	}
	return catalogToResponses(matches, nil), nil
}

func installedSkillNames(skills []string) map[string]struct{} {
	installed := make(map[string]struct{}, len(skills))
	for _, name := range skills {
		installed[strings.ToLower(strings.TrimSpace(name))] = struct{}{}
	}
	return installed
}
