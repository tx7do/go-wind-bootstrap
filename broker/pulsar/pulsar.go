// Package pulsar provides a bootstrap broker builder for Apache Pulsar.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/broker/pulsar"
package pulsar

import (
	"context"
	"fmt"

	pulsarPlugin "github.com/tx7do/go-wind-plugins/transport/pulsar"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterBrokerBuilder(bootstrap.BrokerTypePulsar, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Broker) (any, func(), error) {
	c := cfg.GetPulsar()
	if c == nil {
		return nil, nil, fmt.Errorf("pulsar: config is nil")
	}

	var opts []pulsarPlugin.ServerOption

	if url := c.GetUrl(); url != "" {
		opts = append(opts, pulsarPlugin.WithAddress([]string{url}))
	}

	srv := pulsarPlugin.NewServer(opts...)
	return srv, func() {}, nil
}
