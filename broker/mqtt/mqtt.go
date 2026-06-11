// Package mqtt provides a bootstrap broker builder for MQTT.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/broker/mqtt"
package mqtt

import (
	"context"
	"fmt"

	mqttPlugin "github.com/tx7do/go-wind-plugins/transport/mqtt"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterBrokerBuilder(bootstrap.BrokerTypeMQTT, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Broker) (any, func(), error) {
	c := cfg.GetMqtt()
	if c == nil {
		return nil, nil, fmt.Errorf("mqtt: config is nil")
	}

	var opts []mqttPlugin.ServerOption

	if addr := c.GetAddress(); addr != "" {
		opts = append(opts, mqttPlugin.WithAddress([]string{addr}))
	}
	if clientID := c.GetClientId(); clientID != "" {
		opts = append(opts, mqttPlugin.WithClientId(clientID))
	}
	if username := c.GetUsername(); username != "" {
		opts = append(opts, mqttPlugin.WithAuth(username, c.GetPassword()))
	}

	srv := mqttPlugin.NewServer(opts...)
	return srv, func() {}, nil
}
