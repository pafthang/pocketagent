package tools

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pafthang/pocketagent/pkgs/common"
)

// Config controls exec tool providers.
type Config struct {
	SearchProvider  string
	SerperAPIKey    string
	TavilyAPIKey    string
	CodeExecEnabled bool
	CodeExecTimeout time.Duration
	CodeExecMaxOut  int
	MCPServers      []MCPServerConfig
}

// MCPServerConfig defines an MCP server connection.
type MCPServerConfig struct {
	Name      string            `json:"name"`
	Transport string            `json:"transport"` // stdio | http
	Command   string            `json:"command"`
	Args      []string          `json:"args"`
	URL       string            `json:"url"`
	Env       map[string]string `json:"env"`
	Enabled   bool              `json:"enabled"`
}

// DefaultConfig returns safe defaults (duckduckgo search, code_exec off).
func DefaultConfig() Config {
	return Config{
		SearchProvider:  "auto",
		CodeExecEnabled: !common.IsProduction(),
		CodeExecTimeout: 10 * time.Second,
		CodeExecMaxOut:  8192,
	}
}

// LoadFromEnv reads tool settings from environment variables.
func LoadFromEnv() Config {
	cfg := DefaultConfig()

	cfg.SearchProvider = strings.ToLower(strings.TrimSpace(os.Getenv("SEARCH_PROVIDER")))
	cfg.SerperAPIKey = strings.TrimSpace(os.Getenv("SERPER_API_KEY"))
	cfg.TavilyAPIKey = strings.TrimSpace(os.Getenv("TAVILY_API_KEY"))

	if v := strings.TrimSpace(os.Getenv("CODE_EXEC_ENABLED")); v != "" {
		cfg.CodeExecEnabled = v == "1" || strings.EqualFold(v, "true")
	}
	if v := strings.TrimSpace(os.Getenv("CODE_EXEC_TIMEOUT_SEC")); v != "" {
		if sec, err := parseInt(v); err == nil && sec > 0 {
			cfg.CodeExecTimeout = time.Duration(sec) * time.Second
		}
	}

	if raw := strings.TrimSpace(os.Getenv("MCP_SERVERS")); raw != "" {
		_ = json.Unmarshal([]byte(raw), &cfg.MCPServers)
	}

	return cfg
}

func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}
