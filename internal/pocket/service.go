package pocket

import (
	"errors"
	"net/http"

	"github.com/pocketbase/pocketbase/apis"
)

func Run() error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	app, err := buildApp(cfg)
	if err != nil {
		return err
	}

	err = apis.Serve(app, apis.ServeConfig{
		HttpAddr:        cfg.ListenAddr(),
		ShowStartBanner: true,
	})
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}