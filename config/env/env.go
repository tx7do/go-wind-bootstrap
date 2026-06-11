// Package env provides a bootstrap config action for environment variable config source.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/config/env"
package env

import (
	"context"
	"fmt"

	envPlugin "github.com/tx7do/go-wind-plugins/config/env"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterConfigAction(bootstrap.ConfigTypeEnv, newAction)
}

func newAction(ctx context.Context, cfg *v1.Config) (func(), error) {
	c := cfg.GetEnv()
	if c == nil {
		return nil, fmt.Errorf("env: config is nil")
	}

	var opts []envPlugin.Option

	if prefix := c.GetPrefix(); prefix != "" {
		opts = append(opts, envPlugin.WithPrefix(prefix))
	}
	if key := c.GetKey(); key != "" {
		opts = append(opts, envPlugin.WithKey(key))
	}

	_, err := envPlugin.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("env: create source: %w", err)
	}

	return func() {}, nil
}
