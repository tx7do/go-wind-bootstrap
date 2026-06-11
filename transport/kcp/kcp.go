// Package kcp provides a bootstrap server builder for the KCP (UDP-based reliable) transport.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/transport/kcp"
package kcp

import (
	"fmt"

	kcpPlugin "github.com/tx7do/go-wind-plugins/transport/kcp"
	"github.com/tx7do/go-wind/transport"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeKCP, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetKcp()
	if c == nil {
		return nil, fmt.Errorf("kcp: config is nil")
	}

	addr := c.GetAddr()
	if addr == "" {
		addr = ":9000"
	}

	var opts []kcpPlugin.Option
	opts = append(opts, kcpPlugin.WithAddress(addr))

	if pwd := c.GetBlockCryptPassword(); pwd != "" {
		salt := c.GetBlockCryptSalt()
		opts = append(opts, kcpPlugin.WithBlockCrypt(pwd, salt))
	}
	if ds := c.GetDataShards(); ds > 0 {
		opts = append(opts, kcpPlugin.WithDataShards(int(ds)))
	}
	if ps := c.GetParityShards(); ps > 0 {
		opts = append(opts, kcpPlugin.WithParityShards(int(ps)))
	}

	srv := kcpPlugin.NewServer(opts...)
	return srv, nil
}
