package memo

import (
	"github.com/pafthang/pocketagent/internal/memo/store"
)

// ServiceDeps holds runtime dependencies for the memo HTTP service.
type ServiceDeps struct {
	Config  *Config
	Manager *store.Manager
}

func buildDeps(cfg *Config) (*ServiceDeps, error) {
	mgr, err := store.Open(cfg.DataDir, cfg.Collection, cfg.PersistCompress, cfg.RAGMinSimilarity)
	if err != nil {
		return nil, err
	}
	return &ServiceDeps{Config: cfg, Manager: mgr}, nil
}