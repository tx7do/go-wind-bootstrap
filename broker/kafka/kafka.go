// Package kafka provides a bootstrap broker builder for Apache Kafka.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/broker/kafka"
package kafka

import (
	"context"
	"fmt"

	kafkaPlugin "github.com/tx7do/go-wind-plugins/transport/kafka"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterBrokerBuilder(bootstrap.BrokerTypeKafka, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Broker) (any, func(), error) {
	c := cfg.GetKafka()
	if c == nil {
		return nil, nil, fmt.Errorf("kafka: config is nil")
	}

	var opts []kafkaPlugin.ServerOption

	if len(c.GetBrokers()) > 0 {
		opts = append(opts, kafkaPlugin.WithAddress(c.GetBrokers()))
	}

	switch c.GetAuthType() {
	case "plain":
		opts = append(opts, kafkaPlugin.WithPlainMechanism(c.GetUsername(), c.GetPassword()))
	case "scram_sha256", "scram_sha512":
		opts = append(opts, kafkaPlugin.WithScramMechanism(c.GetAuthType(), c.GetUsername(), c.GetPassword()))
	}

	srv := kafkaPlugin.NewServer(opts...)
	return srv, func() {}, nil
}
