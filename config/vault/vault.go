// Package vault provides a bootstrap config action for HashiCorp Vault config source.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/config/vault"
package vault

import (
	"context"
	"fmt"
	"time"

	vaultapi "github.com/hashicorp/vault/api"

	vaultPlugin "github.com/tx7do/go-wind-plugins/config/vault"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterConfigAction(bootstrap.ConfigTypeVault, newAction)
}

func newAction(ctx context.Context, cfg *v1.Config) (func(), error) {
	c := cfg.GetVault()
	if c == nil {
		return nil, fmt.Errorf("vault: config is nil")
	}

	vaultCfg := vaultapi.DefaultConfig()
	if addr := c.GetAddress(); addr != "" {
		vaultCfg.Address = addr
	}

	client, err := vaultapi.NewClient(vaultCfg)
	if err != nil {
		return nil, fmt.Errorf("vault: create client: %w", err)
	}
	if token := c.GetToken(); token != "" {
		client.SetToken(token)
	}

	var opts []vaultPlugin.Option
	if path := c.GetPath(); path != "" {
		opts = append(opts, vaultPlugin.WithPath(path))
	}
	if dataKey := c.GetDataKey(); dataKey != "" {
		opts = append(opts, vaultPlugin.WithDataKey(dataKey))
	}
	if pollMs := c.GetPollInterval(); pollMs > 0 {
		opts = append(opts, vaultPlugin.WithPollInterval(time.Duration(pollMs)*time.Millisecond))
	}

	_, err = vaultPlugin.New(client, opts...)
	if err != nil {
		return nil, fmt.Errorf("vault: create source: %w", err)
	}

	return func() {}, nil
}
