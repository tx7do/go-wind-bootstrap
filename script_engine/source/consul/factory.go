// Package consul provides Consul KV as a script source.
package consul

import (
	"context"
	"fmt"

	"github.com/tx7do/go-scripts/source"
	consulSource "github.com/tx7do/go-scripts/source/consul"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	scriptEngine "github.com/tx7do/go-wind-bootstrap/script_engine"
)

func init() {
	scriptEngine.MustRegisterSourceFactory(v1.Script_Source_CONSUL, NewSource)
}

// NewSource 根据配置创建 Consul 脚本来源。
func NewSource(cfg *v1.Script_Source) (source.Reader, error) {
	if cfg == nil {
		return nil, fmt.Errorf("consul source: config is nil")
	}

	opts := cfg.GetConsulOptions()
	if opts == nil {
		return nil, fmt.Errorf("consul source: consul_options is required")
	}

	var consulOpts []consulSource.Option

	if addr := opts.GetAddress(); addr != "" {
		consulOpts = append(consulOpts, consulSource.WithAddress(addr))
	}
	if prefix := opts.GetPrefix(); prefix != "" {
		consulOpts = append(consulOpts, consulSource.WithPrefix(prefix))
	}

	return consulSource.New(context.Background(), consulOpts...)
}
