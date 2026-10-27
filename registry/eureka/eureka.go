// Package eureka provides a bootstrap registry action for Eureka service registry.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/registry/eureka"
package eureka

import (
	"context"
	"fmt"

	eurekaPlugin "github.com/tx7do/go-wind-plugins/registry/eureka"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterRegistryAction(bootstrap.RegistryTypeEureka, newAction)
}

func newAction(ctx context.Context, cfg *v1.Registry) (bootstrap.Registry, func(), error) {
	c := cfg.GetEureka()
	if c == nil {
		return nil, nil, fmt.Errorf("eureka: config is nil")
	}

	eps := c.GetEndpoints()
	if len(eps) == 0 {
		return nil, nil, fmt.Errorf("eureka: no endpoints")
	}

	// 选项字段（heartbeat_interval/refresh_interval/eureka_path/max_retry）
	// 由 MapEurekaOptions 生成映射。
	reg, err := eurekaPlugin.New(eps, MapEurekaOptions(c, nil)...)
	if err != nil {
		return nil, nil, fmt.Errorf("eureka: create registry: %w", err)
	}

	return reg, nil, nil
}
