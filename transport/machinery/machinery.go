// Package machinery provides a bootstrap server builder for Machinery distributed task queue.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/transport/machinery"
package machinery

import (
	"fmt"

	machineryPlugin "github.com/tx7do/go-wind-plugins/transport/machinery"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	"github.com/tx7do/go-wind/transport"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeMachinery, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetMachinery()
	if c == nil {
		return nil, fmt.Errorf("machinery: config is nil")
	}

	var opts []machineryPlugin.ServerOption

	brokerType := machineryPlugin.BrokerTypeRedis
	if t := c.GetBrokerType(); t != "" {
		switch t {
		case "amqp":
			brokerType = machineryPlugin.BrokerTypeAmqp
		case "sqs":
			brokerType = machineryPlugin.BrokerTypeSQS
		}
	}

	if addr := c.GetBrokerAddress(); addr != "" {
		opts = append(opts, machineryPlugin.WithBrokerAddress(addr, int(c.GetDb()), brokerType))
	}
	if addr := c.GetResultBackendAddress(); addr != "" {
		opts = append(opts, machineryPlugin.WithResultBackendAddress(addr, int(c.GetDb()), machineryPlugin.BackendTypeRedis))
	}

	srv := machineryPlugin.NewServer(opts...)
	return srv, nil
}
