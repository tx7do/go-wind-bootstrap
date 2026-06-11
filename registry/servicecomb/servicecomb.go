// Package servicecomb provides a bootstrap registry action for ServiceComb service registry.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/registry/servicecomb"
package servicecomb

import (
	"context"
	"fmt"

	"github.com/go-chassis/sc-client"

	servicecombPlugin "github.com/tx7do/go-wind-plugins/registry/servicecomb"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterRegistryAction(bootstrap.RegistryTypeServiceComb, newAction)
}

func newAction(ctx context.Context, appCfg *v1.App, endpoints []string, cfg *v1.Registry) (func(), error) {
	c := cfg.GetServiceComb()
	if c == nil {
		return nil, fmt.Errorf("servicecomb: config is nil")
	}

	eps := c.GetEndpoints()
	if len(eps) == 0 {
		return nil, fmt.Errorf("servicecomb: no endpoints")
	}

	clientOpts := sc.Options{
		Endpoints: eps,
	}

	client, err := sc.NewClient(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("servicecomb: create client: %w", err)
	}

	reg := servicecombPlugin.New(client)
	_ = reg

	return func() {}, nil
}
