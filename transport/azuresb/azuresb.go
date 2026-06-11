// Package azuresb provides a bootstrap server builder for Azure Service Bus transport.
package azuresb

import (
	"fmt"

	azuresbPlugin "github.com/tx7do/go-wind-plugins/transport/azuresb"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	"github.com/tx7do/go-wind/transport"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeAzureSB, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetAzuresb()
	if c == nil {
		return nil, fmt.Errorf("azuresb: config is nil")
	}

	var opts []azuresbPlugin.ServerOption
	if connStr := c.GetConnectionString(); connStr != "" {
		opts = append(opts, azuresbPlugin.WithConnectionString(connStr))
	}
	if codec := c.GetCodec(); codec != "" {
		opts = append(opts, azuresbPlugin.WithCodec(codec))
	}

	return azuresbPlugin.NewServer(opts...), nil
}
