// Package polaris provides a bootstrap registry action for Polaris service registry.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/registry/polaris"
package polaris

import (
	"context"
	"fmt"
	"time"

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

	var polarisCfg config.Configuration
	if configFile := c.GetConfigFile(); configFile != "" {
		loaded, err := config.LoadConfigurationByFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("polaris: load config file: %w", err)
		}
		polarisCfg = loaded
	} else {
		var addresses []string
		if addr := c.GetAddress(); addr != "" {
			addresses = []string{addr}
		}
		polarisCfg = config.NewDefaultConfiguration(addresses)
	}

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
	if priority := c.GetPriority(); priority > 0 {
		opts = append(opts, polarisPlugin.WithPriority(int(priority)))
	}
	if c.Healthy != nil {
		opts = append(opts, polarisPlugin.WithHealthy(c.GetHealthy()))
	}
	if c.Isolate != nil {
		opts = append(opts, polarisPlugin.WithIsolate(c.GetIsolate()))
	}
	if c.Heartbeat != nil {
		opts = append(opts, polarisPlugin.WithHeartbeat(c.GetHeartbeat()))
	}
	if timeout := c.GetTimeout(); timeout > 0 {
		opts = append(opts, polarisPlugin.WithTimeout(time.Duration(timeout)*time.Millisecond))
	}
	if retryCount := c.GetRetryCount(); retryCount > 0 {
		opts = append(opts, polarisPlugin.WithRetryCount(int(retryCount)))
	}
	if ttl := c.GetTtl(); ttl > 0 {
		opts = append(opts, polarisPlugin.WithTTL(int(ttl)))
	} else if c.Heartbeat == nil || c.GetHeartbeat() {
		// 心跳默认开启，TTL 为 0 会导致心跳 ticker panic，兜底一个默认值
		opts = append(opts, polarisPlugin.WithTTL(5))
	}

	reg := polarisPlugin.NewRegistryWithConfig(polarisCfg, opts...)

	regCleanup, err := bootstrap.RegisterInstance(ctx, reg, appCfg, endpoints)
	if err != nil {
		return nil, err
	}
	return regCleanup, nil
}
