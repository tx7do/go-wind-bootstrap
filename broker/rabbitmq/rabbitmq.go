// Package rabbitmq provides a bootstrap broker builder for RabbitMQ.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/broker/rabbitmq"
package rabbitmq

import (
	"context"
	"fmt"

	rabbitmqPlugin "github.com/tx7do/go-wind-plugins/transport/rabbitmq"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterBrokerBuilder(bootstrap.BrokerTypeRabbitMQ, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Broker) (any, func(), error) {
	c := cfg.GetRabbitmq()
	if c == nil {
		return nil, nil, fmt.Errorf("rabbitmq: config is nil")
	}

	var opts []rabbitmqPlugin.ServerOption

	if url := c.GetUrl(); url != "" {
		opts = append(opts, rabbitmqPlugin.WithAddress([]string{url}))
	}
	if exchange := c.GetExchange(); exchange != "" {
		opts = append(opts, rabbitmqPlugin.WithExchange(exchange, true))
	}

	srv := rabbitmqPlugin.NewServer(opts...)
	return srv, func() {}, nil
}
