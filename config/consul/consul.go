// Package consul provides a bootstrap config action for Consul config source.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/config/consul"
package consul

import (
	"context"
	"fmt"

	"github.com/hashicorp/consul/api"

	consulPlugin "github.com/tx7do/go-wind-plugins/config/consul"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterConfigAction(bootstrap.ConfigTypeConsul, newAction)
}

func newAction(ctx context.Context, cfg *v1.Config) (func(), error) {
	c := cfg.GetConsul()
	if c == nil {
		return nil, fmt.Errorf("consul: config is nil")
	}

	consulCfg := api.DefaultConfig()
	if addr := c.GetAddress(); addr != "" {
		consulCfg.Address = addr
	}
	if token := c.GetToken(); token != "" {
		consulCfg.Token = token
	}
	if scheme := c.GetScheme(); scheme != "" {
		consulCfg.Scheme = scheme
	}

	client, err := api.NewClient(consulCfg)
	if err != nil {
		return nil, fmt.Errorf("consul: create client: %w", err)
	}

	var opts []consulPlugin.Option
	if path := c.GetPath(); path != "" {
		opts = append(opts, consulPlugin.WithPath(path))
	}

	_, err = consulPlugin.New(client, opts...)
	if err != nil {
		return nil, fmt.Errorf("consul: create source: %w", err)
	}

	return func() {}, nil
}
