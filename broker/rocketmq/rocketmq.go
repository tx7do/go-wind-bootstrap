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

	opts := MapRocketmqOptions(c, nil)
	// access_key/secret_key/security_token 合成为 WithCredentials 三参选项，非单字段映射，保留手工段。
	if ak := c.GetAccessKey(); ak != "" {
		opts = append(opts, rocketmqPlugin.WithCredentials(ak, c.GetSecretKey(), c.GetSecurityToken()))
	}

	srv := rocketmqPlugin.NewServer(rocketmqOption.DriverTypeV2, opts...)
	return srv, func() {}, nil
}
