// Package consul provides a bootstrap registry action for Consul service registry.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/registry/consul"
package consul

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"

	consulPlugin "github.com/tx7do/go-wind-plugins/registry/consul"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterRegistryAction(bootstrap.RegistryTypeConsul, newAction)
}

func newAction(ctx context.Context, appCfg *v1.App, endpoints []string, cfg *v1.Registry) (func(), error) {
	c := cfg.GetConsul()
	if c == nil {
		return nil, fmt.Errorf("consul: config is nil")
	}

	consulCfg := api.DefaultConfig()
	if addr := c.GetAddress(); addr != "" {
		consulCfg.Address = addr
	}
	if scheme := c.GetScheme(); scheme != "" {
		consulCfg.Scheme = scheme
	}
	if token := c.GetToken(); token != "" {
		consulCfg.Token = token
	}

	client, err := api.NewClient(consulCfg)
	if err != nil {
		return nil, fmt.Errorf("consul: create client: %w", err)
	}

	var opts []consulPlugin.Option
	if c.GetEnableHealthCheck() {
		opts = append(opts, consulPlugin.WithHealthCheck(true))
	}
	if interval := c.GetHealthCheckInterval(); interval > 0 {
		opts = append(opts, consulPlugin.WithHealthCheckInterval(int(interval)))
	}
	if timeout := c.GetHealthCheckTimeout(); timeout > 0 {
		opts = append(opts, consulPlugin.WithTimeout(time.Duration(timeout)*time.Second))
	}
	if c.GetHeartbeat() {
		opts = append(opts, consulPlugin.WithHeartbeat(true))
	}
	if dc := c.GetDatacenter(); dc != "" {
		opts = append(opts, consulPlugin.WithDatacenter(consulPlugin.Datacenter(dc)))
	}
	if dereg := c.GetDeregisterCriticalServiceAfter(); dereg > 0 {
		opts = append(opts, consulPlugin.WithDeregisterCriticalServiceAfter(int(dereg)))
	}

	reg := consulPlugin.New(client, opts...)
	_ = reg

	return func() {}, nil
}
