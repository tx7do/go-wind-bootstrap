// Package hptimer provides a bootstrap server builder for High-Precision Timer.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/transport/hptimer"
package hptimer

import (
	"fmt"

	hptimerPlugin "github.com/tx7do/go-wind-plugins/transport/hptimer"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	"github.com/tx7do/go-wind/transport"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeHPTimer, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetHptimer()
	if c == nil {
		return nil, fmt.Errorf("hptimer: config is nil")
	}

	var opts []hptimerPlugin.Option

	if c.GetGracefullyShutdown() {
		opts = append(opts, hptimerPlugin.WithGracefullyShutdown(true))
	}

	srv := hptimerPlugin.NewServer(opts...)
	return srv, nil
}
