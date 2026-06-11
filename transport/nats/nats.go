// Package nats provides a bootstrap server builder for NATS transport.
package nats

import (
	"fmt"

	natsPlugin "github.com/tx7do/go-wind-plugins/transport/nats"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	"github.com/tx7do/go-wind/transport"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeNATS, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetNats()
	if c == nil {
		return nil, fmt.Errorf("nats: config is nil")
	}

	var opts []natsPlugin.ServerOption
	if addrs := c.GetAddrs(); len(addrs) > 0 {
		opts = append(opts, natsPlugin.WithAddress(addrs))
	}
	if codec := c.GetCodec(); codec != "" {
		opts = append(opts, natsPlugin.WithCodec(codec))
	}

	return natsPlugin.NewServer(opts...), nil
}
