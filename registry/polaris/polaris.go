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

func newAction(ctx context.Context, cfg *v1.Registry) (bootstrap.Registry, func(), error) {
	c := cfg.GetPolaris()
	if c == nil {
		return nil, nil, fmt.Errorf("polaris: config is nil")
	}

	var polarisCfg config.Configuration
	if configFile := c.GetConfigFile(); configFile != "" {
		loaded, err := config.LoadConfigurationByFile(configFile)
		if err != nil {
			return nil, nil, fmt.Errorf("polaris: load config file: %w", err)
		}
		polarisCfg = loaded
	} else {
		var addresses []string
		if addr := c.GetAddress(); addr != "" {
			addresses = []string{addr}
		}
		polarisCfg = config.NewDefaultConfiguration(addresses)
	}

	// 常规选项字段（namespace/token/…/retry_count）由 MapPolarisOptions 生成映射。
	// healthy/isolate/heartbeat 依赖 optional 存在性判断，ttl 有心跳兜底分支，
	// 保留手工段。
	var opts []polarisPlugin.Option
	opts = MapPolarisOptions(c, opts)
	if c.Healthy != nil {
		opts = append(opts, polarisPlugin.WithHealthy(c.GetHealthy()))
	}
	if c.Isolate != nil {
		opts = append(opts, polarisPlugin.WithIsolate(c.GetIsolate()))
	}
	if c.Heartbeat != nil {
		opts = append(opts, polarisPlugin.WithHeartbeat(c.GetHeartbeat()))
	}
	if ttl := c.GetTtl(); ttl > 0 {
		opts = append(opts, polarisPlugin.WithTTL(int(ttl)))
	} else if c.Heartbeat == nil || c.GetHeartbeat() {
		// 心跳默认开启，TTL 为 0 会导致心跳 ticker panic，兜底一个默认值
		opts = append(opts, polarisPlugin.WithTTL(5))
	}

	reg := polarisPlugin.NewRegistryWithConfig(polarisCfg, opts...)

	return reg, nil, nil
}
