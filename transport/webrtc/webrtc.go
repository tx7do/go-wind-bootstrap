// Package webrtc provides a bootstrap server builder for WebRTC.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/transport/webrtc"
package webrtc

import (
	"crypto/tls"
	"fmt"

	webrtcPlugin "github.com/tx7do/go-wind-plugins/transport/webrtc"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	"github.com/tx7do/go-wind/transport"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeWebRTC, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetWebrtc()
	if c == nil {
		return nil, fmt.Errorf("webrtc: config is nil")
	}

	var opts []webrtcPlugin.ServerOption

	if addr := c.GetAddr(); addr != "" {
		opts = append(opts, webrtcPlugin.WithAddress(addr))
	}
	if network := c.GetNetwork(); network != "" {
		opts = append(opts, webrtcPlugin.WithNetwork(network))
	}
	if cert := c.GetTls(); cert != nil {
		opts = append(opts, webrtcPlugin.WithTLSConfig(&tls.Config{}))
	}
	if path := c.GetPath(); path != "" {
		opts = append(opts, webrtcPlugin.WithPath(path))
	}
	if codec := c.GetCodec(); codec != "" {
		opts = append(opts, webrtcPlugin.WithCodec(codec))
	}

	srv := webrtcPlugin.NewServer(opts...)
	return srv, nil
}
