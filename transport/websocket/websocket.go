// Package websocket provides a bootstrap server builder for the WebSocket transport.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/transport/websocket"
package websocket

import (
	"crypto/tls"
	"fmt"

	wsPlugin "github.com/tx7do/go-wind-plugins/transport/websocket"
	"github.com/tx7do/go-wind/transport"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeWebSocket, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetWebsocket()
	if c == nil {
		return nil, fmt.Errorf("websocket: config is nil")
	}

	addr := c.GetAddr()
	if addr == "" {
		addr = ":8080"
	}

	var opts []wsPlugin.Option
	if path := c.GetPath(); path != "" {
		opts = append(opts, wsPlugin.WithPath(path))
	}
	if tlsCfg := buildTLS(c.GetTls()); tlsCfg != nil {
		opts = append(opts, wsPlugin.WithTLSConfig(tlsCfg))
	}

	srv := wsPlugin.NewServer(addr, opts...)
	return srv, nil
}

func buildTLS(tlsCfg *v1.Server_TLS) *tls.Config {
	if tlsCfg == nil {
		return nil
	}
	if f := tlsCfg.GetFile(); f != nil {
		cert, err := tls.LoadX509KeyPair(f.GetCertPath(), f.GetKeyPath())
		if err != nil {
			return nil
		}
		return &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: tlsCfg.GetInsecureSkipVerify(),
		}
	}
	if c := tlsCfg.GetConfig(); c != nil {
		cert, err := tls.X509KeyPair(c.GetCertPem(), c.GetKeyPem())
		if err != nil {
			return nil
		}
		return &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: tlsCfg.GetInsecureSkipVerify(),
		}
	}
	return nil
}
