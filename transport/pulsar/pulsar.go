// Package pulsar provides a bootstrap server builder for Pulsar transport.
package pulsar

import (
	"fmt"

	pulsarPlugin "github.com/tx7do/go-wind-plugins/transport/pulsar"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	"github.com/tx7do/go-wind/transport"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypePulsar, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetPulsar()
	if c == nil {
		return nil, fmt.Errorf("pulsar: config is nil")
	}

	var opts []pulsarPlugin.ServerOption
	if addrs := c.GetAddrs(); len(addrs) > 0 {
		opts = append(opts, pulsarPlugin.WithAddress(addrs))
	}
	if codec := c.GetCodec(); codec != "" {
		opts = append(opts, pulsarPlugin.WithCodec(codec))
	}

	return pulsarPlugin.NewServer(opts...), nil
}
