// Package http3 provides a bootstrap server builder for the HTTP/3 (QUIC) transport.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/transport/http3"
package http3

import (
	"crypto/tls"
	"fmt"

	http3Plugin "github.com/tx7do/go-wind-plugins/transport/http3"
	"github.com/tx7do/go-wind/transport"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeHTTP3, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetHttp3()
	if c == nil {
		return nil, fmt.Errorf("http3: config is nil")
	}

	addr := c.GetAddr()
	if addr == "" {
		addr = ":443"
	}

	var opts []http3Plugin.ServerOption
	opts = append(opts, http3Plugin.WithAddress(addr))

	if tlsCfg := buildTLS(c.GetTls()); tlsCfg != nil {
		opts = append(opts, http3Plugin.WithTLSConfig(tlsCfg))
	}

	srv := http3Plugin.NewServer(opts...)
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
