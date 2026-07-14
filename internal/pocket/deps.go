package pocket

import (
	"log"

	"github.com/pafthang/pocketagent/internal/pocket/bootstrap"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func buildApp(cfg *Config) (*pocketbase.PocketBase, error) {
	dataDir := cfg.DataDir
	if dataDir == "" {
		dataDir = "data"
	}

	app := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDataDir: dataDir,
	})

	app.OnBootstrap().BindFunc(func(e *core.BootstrapEvent) error {
		if err := e.Next(); err != nil {
			return err
		}
		return bootstrap.Run(e.App, bootstrap.Config{
			SuperuserEmail:    cfg.SuperuserEmail,
			SuperuserPassword: cfg.SuperuserPassword,
		})
	})

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		log.Printf("PocketBase starting on %s (data dir: %s)", cfg.ListenAddr(), dataDir)
		return e.Next()
	})

	if err := app.Bootstrap(); err != nil {
		return nil, err
	}

	return app, nil
}