package bootstrap

import (
	"fmt"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func resolveTracer(cfg *v1.Tracer) (interface{}, func(), error) {
	if cfg.GetType() == "" {
		return nil, nil, fmt.Errorf("bootstrap: tracer type not specified")
	}

	b, err := getTracerBuilder(cfg.GetType())
	if err != nil {
		return nil, nil, err
	}

	tp, cleanup, err := b(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("bootstrap: build tracer %q: %w", cfg.GetType(), err)
	}
	return tp, cleanup, nil
}
