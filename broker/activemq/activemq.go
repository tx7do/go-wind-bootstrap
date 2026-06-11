// Package activemq provides a bootstrap broker builder for Apache ActiveMQ.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/broker/activemq"
package activemq

import (
	"context"
	"fmt"

	activemqPlugin "github.com/tx7do/go-wind-plugins/transport/activemq"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterBrokerBuilder(bootstrap.BrokerTypeActiveMQ, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Broker) (any, func(), error) {
	c := cfg.GetActivemq()
	if c == nil {
		return nil, nil, fmt.Errorf("activemq: config is nil")
	}

	var opts []activemqPlugin.ServerOption

	if addr := c.GetAddress(); addr != "" {
		opts = append(opts, activemqPlugin.WithAddress([]string{addr}))
	}

	srv := activemqPlugin.NewServer(opts...)
	return srv, func() {}, nil
}
