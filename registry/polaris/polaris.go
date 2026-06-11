// Package polaris provides a bootstrap registry action for Polaris service registry.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/registry/polaris"
package polaris

import (
	"context"
	"fmt"

	"github.com/polarismesh/polaris-go/pkg/config"

	polarisPlugin "github.com/tx7do/go-wind-plugins/registry/polaris"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterRegistryAction(bootstrap.RegistryTypePolaris, newAction)
}

func newAction(ctx context.Context, appCfg *v1.App, endpoints []string, cfg *v1.Registry) (func(), error) {
	c := cfg.GetPolaris()
	if c == nil {
		return nil, fmt.Errorf("polaris: config is nil")
	}

	var addresses []string
	if addr := c.GetAddress(); addr != "" {
		addresses = []string{addr}
	}

	polarisCfg := config.NewDefaultConfiguration(addresses)

	var opts []polarisPlugin.Option
	if ns := c.GetNamespace(); ns != "" {
		opts = append(opts, polarisPlugin.WithNamespace(ns))
	}
	if token := c.GetToken(); token != "" {
		opts = append(opts, polarisPlugin.WithServiceToken(token))
	}
	if protocol := c.GetProtocol(); protocol != "" {
		opts = append(opts, polarisPlugin.WithProtocol(protocol))
	}
	if weight := c.GetWeight(); weight > 0 {
		opts = append(opts, polarisPlugin.WithWeight(int(weight)))
	}

	reg := polarisPlugin.NewRegistryWithConfig(polarisCfg, opts...)
	_ = reg

	return func() {}, nil
}
