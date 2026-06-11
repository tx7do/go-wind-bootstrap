// Package mqtt provides a bootstrap server builder for MQTT transport.
package mqtt

import (
	"fmt"

	mqttPlugin "github.com/tx7do/go-wind-plugins/transport/mqtt"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	"github.com/tx7do/go-wind/transport"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeMQTT, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetMqtt()
	if c == nil {
		return nil, fmt.Errorf("mqtt: config is nil")
	}

	var opts []mqttPlugin.ServerOption
	if addrs := c.GetAddrs(); len(addrs) > 0 {
		opts = append(opts, mqttPlugin.WithAddress(addrs))
	}
	if codec := c.GetCodec(); codec != "" {
		opts = append(opts, mqttPlugin.WithCodec(codec))
	}

	return mqttPlugin.NewServer(opts...), nil
}
