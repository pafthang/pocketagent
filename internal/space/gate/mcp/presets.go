package mcpapis

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/pafthang/pocketagent/pkgs/common"
)

type mcpPresetResponse struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Icon        string            `json:"icon"`
	Category    string            `json:"category"`
	Package     string            `json:"package"`
	Transport   string            `json:"transport"`
	URL         string            `json:"url,omitempty"`
	DocsURL     string            `json:"docs_url"`
	NeedsArgs   bool              `json:"needs_args"`
	OAuth       bool              `json:"oauth"`
	Installed   bool              `json:"installed"`
	EnvKeys     []mcpPresetEnvKey `json:"env_keys"`
}

type mcpPresetEnvKey struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Required    bool   `json:"required"`
	Placeholder string `json:"placeholder"`
	Secret      bool   `json:"secret"`
}

type mcpPresetFile struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Icon        string            `json:"icon"`
	Category    string            `json:"category"`
	Package     string            `json:"package"`
	Transport   string            `json:"transport"`
	URL         string            `json:"url,omitempty"`
	DocsURL     string            `json:"docs_url"`
	NeedsArgs   bool              `json:"needs_args"`
	OAuth       bool              `json:"oauth"`
	EnvKeys     []mcpPresetEnvKey `json:"env_keys"`
}

func loadMCPPresents(installed map[string]struct{}) ([]mcpPresetResponse, error) {
	dir, err := common.FindConfigsDir()
	if err != nil {
		return nil, err
	}

	raw, err := os.ReadFile(filepath.Join(dir, "mcp-presets.json"))
	if err != nil {
		return nil, err
	}

	var presets []mcpPresetFile
	if err := json.Unmarshal(raw, &presets); err != nil {
		return nil, err
	}

	result := make([]mcpPresetResponse, 0, len(presets))
	for _, preset := range presets {
		_, isInstalled := installed[strings.ToLower(strings.TrimSpace(preset.Name))]
		result = append(result, mcpPresetResponse{
			ID:          preset.ID,
			Name:        preset.Name,
			Description: preset.Description,
			Icon:        preset.Icon,
			Category:    preset.Category,
			Package:     preset.Package,
			Transport:   preset.Transport,
			URL:         preset.URL,
			DocsURL:     preset.DocsURL,
			NeedsArgs:   preset.NeedsArgs,
			OAuth:       preset.OAuth,
			Installed:   isInstalled,
			EnvKeys:     preset.EnvKeys,
		})
	}
	return result, nil
}

func installedServerNames(servers []string) map[string]struct{} {
	installed := make(map[string]struct{}, len(servers))
	for _, name := range servers {
		installed[strings.ToLower(strings.TrimSpace(name))] = struct{}{}
	}
	return installed
}
