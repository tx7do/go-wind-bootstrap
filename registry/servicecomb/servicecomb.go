// Package servicecomb provides a bootstrap registry action for ServiceComb service registry.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/registry/servicecomb"
package servicecomb

import (
	"context"
	"fmt"
	"time"

	"github.com/go-chassis/cari/rbac"
	"github.com/go-chassis/sc-client"

	servicecombPlugin "github.com/tx7do/go-wind-plugins/registry/servicecomb"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	"github.com/tx7do/go-wind-bootstrap/tlsutil"
)

func init() {
	bootstrap.MustRegisterRegistryAction(bootstrap.RegistryTypeServiceComb, newAction)
}

func newAction(ctx context.Context, cfg *v1.Registry) (bootstrap.Registry, func(), error) {
	c := cfg.GetServiceComb()
	if c == nil {
		return nil, nil, fmt.Errorf("servicecomb: config is nil")
	}

	eps := c.GetEndpoints()
	if len(eps) == 0 {
		return nil, nil, fmt.Errorf("servicecomb: no endpoints")
	}

	clientOpts := sc.Options{
		Endpoints:  eps,
		EnableSSL:  c.GetEnableSsl(),
		Compressed: c.GetCompressed(),
		Verbose:    c.GetVerbose(),
	}
	if timeout := c.GetTimeout(); timeout > 0 {
		clientOpts.Timeout = time.Duration(timeout) * time.Second
	}
	if c.GetEnableAuth() {
		clientOpts.EnableAuth = true
		clientOpts.AuthUser = &rbac.AuthUser{
			Username: c.GetUsername(),
			Password: c.GetPassword(),
		}
	}
	if expiration := c.GetTokenExpiration(); expiration > 0 {
		clientOpts.TokenExpiration = time.Duration(expiration) * time.Second
	}
	if tlsCfg, err := tlsutil.ClientTLS(c.GetTls()); err != nil {
		return nil, nil, fmt.Errorf("servicecomb: load tls: %w", err)
	} else if tlsCfg != nil {
		clientOpts.TLSConfig = tlsCfg
	}

	client, err := sc.NewClient(clientOpts)
	if err != nil {
		return nil, nil, fmt.Errorf("servicecomb: create client: %w", err)
	}

	reg := servicecombPlugin.New(client)

	return reg, nil, nil
}
