// Package eureka provides a bootstrap registry action for Eureka service registry.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/registry/eureka"
package eureka

import (
	"context"
	"fmt"
	"time"

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

	var opts []eurekaPlugin.Option
	if interval := c.GetHeartbeatInterval(); interval > 0 {
		opts = append(opts, eurekaPlugin.WithHeartbeat(time.Duration(interval)*time.Second))
	}
	if ri := c.GetRefreshInterval(); ri > 0 {
		opts = append(opts, eurekaPlugin.WithRefresh(time.Duration(ri)*time.Second))
	}
	if path := c.GetEurekaPath(); path != "" {
		opts = append(opts, eurekaPlugin.WithEurekaPath(path))
	}
	if mr := c.GetMaxRetry(); mr > 0 {
		opts = append(opts, eurekaPlugin.MaxRetry(int(mr)))
	}

	reg, err := eurekaPlugin.New(eps, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("eureka: create registry: %w", err)
	}

	return reg, nil, nil
}
