// Package trpc provides a bootstrap server builder for the tRPC (Tencent) transport.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/transport/trpc"
package trpc

import (
	"fmt"

	trpcPlugin "github.com/tx7do/go-wind-plugins/transport/trpc"
	"github.com/tx7do/go-wind/transport"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeTRPC, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetTrpc()
	if c == nil {
		return nil, fmt.Errorf("trpc: config is nil")
	}

	addr := c.GetAddr()
	if addr == "" {
		addr = ":8000"
	}

	var opts []trpcPlugin.ServerOption
	opts = append(opts, trpcPlugin.WithAddress(addr))

	srv := trpcPlugin.NewServer(opts...)
	return srv, nil
}
