// Package activemq provides a bootstrap server builder for ActiveMQ transport.
package activemq

import (
	"fmt"

	activemqPlugin "github.com/tx7do/go-wind-plugins/transport/activemq"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	"github.com/tx7do/go-wind/transport"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeActiveMQ, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetActivemq()
	if c == nil {
		return nil, fmt.Errorf("activemq: config is nil")
	}

	var opts []activemqPlugin.ServerOption
	if addrs := c.GetAddrs(); len(addrs) > 0 {
		opts = append(opts, activemqPlugin.WithAddress(addrs))
	}
	if codec := c.GetCodec(); codec != "" {
		opts = append(opts, activemqPlugin.WithCodec(codec))
	}

	return activemqPlugin.NewServer(opts...), nil
}
