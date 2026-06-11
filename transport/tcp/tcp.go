// Package tcp provides a bootstrap server builder for the TCP socket transport.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/transport/tcp"
package tcp

import (
	"crypto/tls"
	"fmt"

	tcpPlugin "github.com/tx7do/go-wind-plugins/transport/tcp"
	"github.com/tx7do/go-wind/transport"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeTCP, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetTcp()
	if c == nil {
		return nil, fmt.Errorf("tcp: config is nil")
	}

	addr := c.GetAddr()
	if addr == "" {
		addr = ":9000"
	}

	var opts []tcpPlugin.Option
	opts = append(opts, tcpPlugin.WithAddress(addr))

	if tlsCfg := buildTLS(c.GetTls()); tlsCfg != nil {
		opts = append(opts, tcpPlugin.WithTLSConfig(tlsCfg))
	}

	srv := tcpPlugin.NewServer(opts...)
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
