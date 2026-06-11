// Package polaris provides a bootstrap config action for Polaris config source.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/config/polaris"
package polaris

import (
	"context"
	"fmt"

	"github.com/polarismesh/polaris-go"
	"github.com/polarismesh/polaris-go/pkg/config"

	polarisPlugin "github.com/tx7do/go-wind-plugins/config/polaris"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterConfigAction(bootstrap.ConfigTypePolaris, newAction)
}

func newAction(ctx context.Context, cfg *v1.Config) (func(), error) {
	c := cfg.GetPolaris()
	if c == nil {
		return nil, fmt.Errorf("polaris: config is nil")
	}

	var addresses []string
	if addr := c.GetAddress(); addr != "" {
		addresses = []string{addr}
	}

	polarisCfg := config.NewDefaultConfiguration(addresses)

	sdk, err := polaris.NewConfigAPIByConfig(polarisCfg)
	if err != nil {
		return nil, fmt.Errorf("polaris: create SDK: %w", err)
	}

	var opts []polarisPlugin.Option
	if ns := c.GetNamespace(); ns != "" {
		opts = append(opts, polarisPlugin.WithNamespace(ns))
	}
	if fg := c.GetGroup(); fg != "" {
		opts = append(opts, polarisPlugin.WithFileGroup(fg))
	}
	if fn := c.GetFileName(); fn != "" {
		opts = append(opts, polarisPlugin.WithFileName(fn))
	}

	_, err = polarisPlugin.New(sdk, opts...)
	if err != nil {
		return nil, fmt.Errorf("polaris: create source: %w", err)
	}

	return func() {}, nil
}
