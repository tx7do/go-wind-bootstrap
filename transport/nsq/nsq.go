// Package nsq provides a bootstrap server builder for NSQ transport.
package nsq

import (
	"fmt"

	nsqPlugin "github.com/tx7do/go-wind-plugins/transport/nsq"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	"github.com/tx7do/go-wind/transport"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeNSQ, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetNsq()
	if c == nil {
		return nil, fmt.Errorf("nsq: config is nil")
	}

	var opts []nsqPlugin.ServerOption
	if addrs := c.GetAddrs(); len(addrs) > 0 {
		opts = append(opts, nsqPlugin.WithAddress(addrs))
	}
	if lookupd := c.GetLookupdAddrs(); len(lookupd) > 0 {
		opts = append(opts, nsqPlugin.WithLookupdAddress(lookupd))
	}
	if codec := c.GetCodec(); codec != "" {
		opts = append(opts, nsqPlugin.WithCodec(codec))
	}

	return nsqPlugin.NewServer(opts...), nil
}
