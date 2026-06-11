package bootstrap

import (
	"context"
	"fmt"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func resolveRegistry(ctx context.Context, cfg *v1.Registry, appCfg *v1.App) (func(), error) {
	if cfg.GetType() == "" {
		return nil, fmt.Errorf("bootstrap: registry type not specified")
	}

	a, err := getRegistryAction(cfg.GetType())
	if err != nil {
		return nil, err
	}

	var endpoints []string
	cleanup, err := a(ctx, appCfg, endpoints)
	if err != nil {
		return nil, fmt.Errorf("bootstrap: build registry %q: %w", cfg.GetType(), err)
	}
	return cleanup, nil
}
