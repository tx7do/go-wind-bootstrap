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
	// 常规连接字段 → 客户端配置的映射由 protoc-gen-fieldmap 生成（mapping.spec.yaml）。
	FillConsulConfig(c, consulCfg)

	// TLS 展开到 Consul 专属 TLSConfig 子结构，保留手工段。
	if tlsCfg := c.GetTls(); tlsCfg != nil {
		consulCfg.TLSConfig = api.TLSConfig{
			InsecureSkipVerify: tlsCfg.GetInsecureSkipVerify(),
		}
		if f := tlsCfg.GetFile(); f != nil {
			consulCfg.TLSConfig.CertFile = f.GetCertPath()
			consulCfg.TLSConfig.KeyFile = f.GetKeyPath()
			consulCfg.TLSConfig.CAFile = f.GetCaPath()
		}
		if cf := tlsCfg.GetConfig(); cf != nil {
			consulCfg.TLSConfig.CertPEM = cf.GetCertPem()
			consulCfg.TLSConfig.KeyPEM = cf.GetKeyPem()
			consulCfg.TLSConfig.CAPem = cf.GetCaPem()
		}
	}

	client, err := api.NewClient(consulCfg)
	if err != nil {
		return nil, fmt.Errorf("consul: create client: %w", err)
	}

	// path 为插件选项，保留手工段。
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
