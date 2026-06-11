// Package nsq provides a bootstrap broker builder for NSQ.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/broker/nsq"
package nsq

import (
	"context"
	"fmt"

	nsqPlugin "github.com/tx7do/go-wind-plugins/transport/nsq"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterBrokerBuilder(bootstrap.BrokerTypeNSQ, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Broker) (any, func(), error) {
	c := cfg.GetNsq()
	if c == nil {
		return nil, nil, fmt.Errorf("nsq: config is nil")
	}

	var opts []nsqPlugin.ServerOption

	if len(c.GetAddrs()) > 0 {
		opts = append(opts, nsqPlugin.WithAddress(c.GetAddrs()))
	}
	if len(c.GetLookupdAddrs()) > 0 {
		opts = append(opts, nsqPlugin.WithLookupdAddress(c.GetLookupdAddrs()))
	}

	srv := nsqPlugin.NewServer(opts...)
	return srv, func() {}, nil
}
