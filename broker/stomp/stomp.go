// Package stomp provides a bootstrap broker builder for STOMP protocol.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/broker/stomp"
package stomp

import (
	"context"
	"fmt"

	"github.com/tx7do/go-wind-plugins/broker"
	stompPlugin "github.com/tx7do/go-wind-plugins/broker/stomp"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterBrokerBuilder(bootstrap.BrokerTypeSTOMP, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Broker) (any, func(), error) {
	c := cfg.GetStomp()
	if c == nil {
		return nil, nil, fmt.Errorf("stomp: config is nil")
	}

	var opts []broker.Option

	if addr := c.GetAddress(); addr != "" {
		opts = append(opts, broker.WithAddress(addr))
	}
	if username := c.GetUsername(); username != "" {
		opts = append(opts, stompPlugin.WithAuth(username, c.GetPassword()))
	}
	if vhost := c.GetVhost(); vhost != "" {
		opts = append(opts, stompPlugin.WithVirtualHost(vhost))
	}

	b := stompPlugin.NewBroker(opts...)
	return b, func() {}, nil
}
