// Package vault provides a bootstrap config action for HashiCorp Vault config source.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/config/vault"
package vault

import (
	"context"
	"fmt"

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

	// address → vaultapi.Config 仅一个字段，保留手工段。
	vaultCfg := vaultapi.DefaultConfig()
	if addr := c.GetAddress(); addr != "" {
		vaultCfg.Address = addr
	}

	client, err := vaultapi.NewClient(vaultCfg)
	if err != nil {
		return nil, fmt.Errorf("vault: create client: %w", err)
	}
	// token 是客户端方法调用，保留手工段。
	if token := c.GetToken(); token != "" {
		client.SetToken(token)
	}

	// 常规配置字段 → 选项的映射由 protoc-gen-fieldmap 生成（mapping.spec.yaml）。
	var opts []vaultPlugin.Option
	opts = MapVaultOptions(c, opts)

	_, err = vaultPlugin.New(client, opts...)
	if err != nil {
		return nil, fmt.Errorf("vault: create source: %w", err)
	}

	return func() {}, nil
}
