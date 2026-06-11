// Package thrift provides a bootstrap server builder for the Thrift RPC transport.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/transport/thrift"
package thrift

import (
	"crypto/tls"
	"fmt"

	thriftPlugin "github.com/tx7do/go-wind-plugins/transport/thrift"
	"github.com/tx7do/go-wind/transport"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeThrift, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetThrift()
	if c == nil {
		return nil, fmt.Errorf("thrift: config is nil")
	}

	addr := c.GetAddr()
	if addr == "" {
		addr = ":7700"
	}

	var opts []thriftPlugin.Option
	if proto := c.GetProtocol(); proto != "" {
		opts = append(opts, thriftPlugin.WithProtocol(proto))
	}
	if c.GetBuffered() || c.GetFramed() || c.GetBufferSize() > 0 {
		bufSize := int(c.GetBufferSize())
		if bufSize == 0 {
			bufSize = 8192
		}
		opts = append(opts, thriftPlugin.WithTransportConfig(c.GetBuffered(), c.GetFramed(), bufSize))
	}
	if tlsCfg := buildTLS(c.GetTls()); tlsCfg != nil {
		opts = append(opts, thriftPlugin.WithTLSConfig(tlsCfg))
	}

	srv := thriftPlugin.NewServer(addr, opts...)
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
