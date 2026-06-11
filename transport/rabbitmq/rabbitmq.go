// Package rabbitmq provides a bootstrap server builder for RabbitMQ transport.
package rabbitmq

import (
	"fmt"

	rabbitmqPlugin "github.com/tx7do/go-wind-plugins/transport/rabbitmq"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	"github.com/tx7do/go-wind/transport"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeRabbitMQ, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetRabbitmq()
	if c == nil {
		return nil, fmt.Errorf("rabbitmq: config is nil")
	}

	var opts []rabbitmqPlugin.ServerOption
	if addrs := c.GetAddrs(); len(addrs) > 0 {
		opts = append(opts, rabbitmqPlugin.WithAddress(addrs))
	}
	if codec := c.GetCodec(); codec != "" {
		opts = append(opts, rabbitmqPlugin.WithCodec(codec))
	}

	return rabbitmqPlugin.NewServer(opts...), nil
}
