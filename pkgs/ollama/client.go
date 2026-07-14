package ollama

import (
	"github.com/pafthang/pocketagent/pkgs/common"
	"github.com/pafthang/pocketagent/pkgs/ollama/api"
)

// Client talks to a local or remote Ollama HTTP API.
type Client = api.Client

// NewClient creates a client with default settings.
func NewClient(url string) *Client { return api.New(url) }

// NewConfigured creates a client with an optional embedding model override.
func NewConfigured(url, embedModel string) *Client { return api.NewConfigured(url, embedModel) }

// CircuitBreaker is re-exported for optional client tuning.
type CircuitBreaker = common.CircuitBreaker