package fileapis

import (
	memoclient "github.com/pafthang/pocketagent/internal/memo/client"
	"github.com/pafthang/pocketagent/internal/files/blob"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/ollama"
)

// Deps holds runtime dependencies for file HTTP handlers.
type Deps struct {
	PB     *pbclient.Client
	Store  blob.Backend
	Memo   *memoclient.Client
	Ollama *ollama.Client
}