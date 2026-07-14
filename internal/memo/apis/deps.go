package memoapis

import (
	memoclient "github.com/pafthang/pocketagent/internal/memo/client"
	"github.com/pafthang/pocketagent/pkgs/ollama"
)

// Deps holds gate-facing memo HTTP handler dependencies.
type Deps struct {
	Memo   *memoclient.Client
	Ollama *ollama.Client
}