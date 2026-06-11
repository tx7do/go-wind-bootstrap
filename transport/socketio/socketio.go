// Package socketio provides a bootstrap server builder for Socket.IO.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/transport/socketio"
package socketio

import (
	"crypto/tls"
	"fmt"

	socketioPlugin "github.com/tx7do/go-wind-plugins/transport/socketio"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	"github.com/tx7do/go-wind/transport"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeSocketIO, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetSocketio()
	if c == nil {
		return nil, fmt.Errorf("socketio: config is nil")
	}

	var opts []socketioPlugin.Option

	if addr := c.GetAddr(); addr != "" {
		opts = append(opts, socketioPlugin.WithAddress(addr))
	}
	if network := c.GetNetwork(); network != "" {
		opts = append(opts, socketioPlugin.WithNetwork(network))
	}
	if cert := c.GetTls(); cert != nil {
		opts = append(opts, socketioPlugin.WithTLSConfig(&tls.Config{}))
	}
	if codec := c.GetCodec(); codec != "" {
		opts = append(opts, socketioPlugin.WithCodec(codec))
	}
	if path := c.GetPath(); path != "" {
		opts = append(opts, socketioPlugin.WithPath(path))
	}

	srv := socketioPlugin.NewServer(opts...)
	return srv, nil
}
