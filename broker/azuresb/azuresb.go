// Package azuresb provides a bootstrap broker builder for Azure Service Bus.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/broker/azuresb"
package azuresb

import (
	"context"
	"fmt"

	azuresbPlugin "github.com/tx7do/go-wind-plugins/transport/azuresb"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterBrokerBuilder(bootstrap.BrokerTypeAzureSB, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Broker) (any, func(), error) {
	c := cfg.GetAzuresb()
	if c == nil {
		return nil, nil, fmt.Errorf("azuresb: config is nil")
	}

	var opts []azuresbPlugin.ServerOption

	if connStr := c.GetConnectionString(); connStr != "" {
		opts = append(opts, azuresbPlugin.WithConnectionString(connStr))
	}

	srv := azuresbPlugin.NewServer(opts...)
	return srv, func() {}, nil
}
