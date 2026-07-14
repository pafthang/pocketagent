package gateapis

import (
	agentclient "github.com/pafthang/pocketagent/internal/agent/client"
	filesclient "github.com/pafthang/pocketagent/internal/files/client"
	memoapis "github.com/pafthang/pocketagent/internal/memo/apis"
	natsclient "github.com/pafthang/pocketagent/internal/nats/client"
	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	spaceclient "github.com/pafthang/pocketagent/internal/space/client"
	"github.com/pafthang/pocketagent/pkgs/common"
	apimw "github.com/pafthang/pocketagent/pkgs/middle"
)

// Deps holds runtime dependencies for gate HTTP routes.
type Deps struct {
	NATS       *natsclient.Client
	PB         *pbclient.Client
	Space      *spaceclient.Client
	Agent      *agentclient.Client
	Files      *filesclient.Client
	RBAC       *apimw.PocketRBAC
	Memo       memoapis.Deps
	EmbedModel string
	OllamaURL  string
	LLMModel   string
	RateLimit  common.RateLimitConfig
}