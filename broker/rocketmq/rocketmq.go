// Package rocketmq provides a bootstrap broker builder for Apache RocketMQ.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/broker/rocketmq"
package rocketmq

import (
	"context"
	"fmt"

	rocketmqOption "github.com/tx7do/go-wind-plugins/broker/rocketmq/option"
	rocketmqPlugin "github.com/tx7do/go-wind-plugins/transport/rocketmq"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterBrokerBuilder(bootstrap.BrokerTypeRocketMQ, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Broker) (any, func(), error) {
	c := cfg.GetRocketmq()
	if c == nil {
		return nil, nil, fmt.Errorf("rocketmq: config is nil")
	}

	var opts []rocketmqPlugin.ServerOption

	if len(c.GetNameServers()) > 0 {
		opts = append(opts, rocketmqPlugin.WithNameServer(c.GetNameServers()))
	}
	if url := c.GetNameServerUrl(); url != "" {
		opts = append(opts, rocketmqPlugin.WithNameServerDomain(url))
	}
	if ak := c.GetAccessKey(); ak != "" {
		opts = append(opts, rocketmqPlugin.WithCredentials(ak, c.GetSecretKey(), c.GetSecurityToken()))
	}
	if ns := c.GetNamespace(); ns != "" {
		opts = append(opts, rocketmqPlugin.WithNamespace(ns))
	}
	if inst := c.GetInstanceName(); inst != "" {
		opts = append(opts, rocketmqPlugin.WithInstanceName(inst))
	}
	if gn := c.GetGroupName(); gn != "" {
		opts = append(opts, rocketmqPlugin.WithGroupName(gn))
	}
	if rc := c.GetRetryCount(); rc > 0 {
		opts = append(opts, rocketmqPlugin.WithRetryCount(int(rc)))
	}

	srv := rocketmqPlugin.NewServer(rocketmqOption.DriverTypeV2, opts...)
	return srv, func() {}, nil
}
