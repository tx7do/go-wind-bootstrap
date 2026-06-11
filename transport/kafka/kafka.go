// Package kafka provides a bootstrap server builder for Kafka transport.
package kafka

import (
	"fmt"

	kafkaPlugin "github.com/tx7do/go-wind-plugins/transport/kafka"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	"github.com/tx7do/go-wind/transport"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeKafka, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetKafka()
	if c == nil {
		return nil, fmt.Errorf("kafka: config is nil")
	}

	var opts []kafkaPlugin.ServerOption
	if addrs := c.GetAddrs(); len(addrs) > 0 {
		opts = append(opts, kafkaPlugin.WithAddress(addrs))
	}
	if codec := c.GetCodec(); codec != "" {
		opts = append(opts, kafkaPlugin.WithCodec(codec))
	}

	return kafkaPlugin.NewServer(opts...), nil
}
