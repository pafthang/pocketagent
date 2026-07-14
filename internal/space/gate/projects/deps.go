package projectapis

import (
	natsclient "github.com/pafthang/pocketagent/internal/nats/client"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/ollama"
)

// Deps wires project HTTP handlers.
type Deps struct {
	PB       *pbclient.Client
	NC       *natsclient.Client
	Ollama   *ollama.Client
	LLMModel string
}
