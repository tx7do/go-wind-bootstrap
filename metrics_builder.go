package bootstrap

import (
	"fmt"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func resolveMetrics(cfg *v1.Metrics) (func(), error) {
	if cfg.GetType() == "" {
		return nil, fmt.Errorf("bootstrap: metrics type not specified")
	}

	b, err := getMetricsBuilder(cfg.GetType())
	if err != nil {
		return nil, err
	}

	cleanup, err := b(cfg)
	if err != nil {
		return nil, fmt.Errorf("bootstrap: build metrics %q: %w", cfg.GetType(), err)
	}
	return cleanup, nil
}
