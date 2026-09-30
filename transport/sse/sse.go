// Package sse provides a bootstrap server builder for the SSE (Server-Sent Events) transport.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/transport/sse"
//
// The declarative config covers the listen address, the subscription path and
// TLS. Two registration hooks exist for everything the config cannot express —
// both must be called before [bootstrap.Bootstrap]:
//
//   - [RegisterServerOption] appends constructor options. Use it for the
//     auth chain (WithTokenExtractor / WithAuthorizeFunc / WithSubscriberFunction /
//     WithUnSubscriberFunction) and stream behaviour switches — these are
//     constructor-time options that cannot be applied to a built server.
//     Registered options are applied AFTER the config-derived ones (last wins).
//
//   - [RegisterServerSetup] appends post-construction callbacks, invoked in
//     registration order once the server is built. Use it to attach HTTP
//     middleware (srv.Use), register plain routes (srv.HandleFunc) and
//     pre-create streams (srv.CreateStream).
package sse

import (
	"crypto/tls"
	"fmt"

	ssePlugin "github.com/tx7do/go-wind-plugins/transport/sse"
	"github.com/tx7do/go-wind/transport"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

// Server is a type alias for the plugins SSE Server.
// User code should reference this type so that it only needs to depend on
// the adapter package.
type Server = ssePlugin.Server

// serverOptions holds constructor options appended by [RegisterServerOption].
var serverOptions []ssePlugin.Option

// RegisterServerOption appends constructor options applied when the SSE server
// is built, after the config-derived options (last wins). Use it for the auth
// chain and stream behaviour switches, which are constructor-time options.
//
// This function must be called before [bootstrap.Bootstrap].
func RegisterServerOption(opts ...ssePlugin.Option) {
	serverOptions = append(serverOptions, opts...)
}

// serverSetups holds callbacks for configuring the SSE server after
// construction. They are called once during server construction, after the
// server is created but before it is returned to the bootstrap framework.
var serverSetups []func(srv *Server)

// RegisterServerSetup appends a post-construction callback, invoked in
// registration order. Use it to attach HTTP middleware (srv.Use), register
// plain routes (srv.HandleFunc) and pre-create streams (srv.CreateStream).
//
// This function must be called before [bootstrap.Bootstrap].
func RegisterServerSetup(fn func(srv *Server)) {
	serverSetups = append(serverSetups, fn)
}

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeSSE, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetSse()
	if c == nil {
		return nil, fmt.Errorf("sse: config is nil")
	}

	addr := c.GetAddr()
	if addr == "" {
		addr = ":8080"
	}

	var opts []ssePlugin.Option
	if path := c.GetPath(); path != "" {
		opts = append(opts, ssePlugin.WithPath(path))
	}
	if tlsCfg := buildTLS(c.GetTls()); tlsCfg != nil {
		opts = append(opts, ssePlugin.WithTLSConfig(tlsCfg))
	}

	// Application-level constructor options come last (last wins).
	opts = append(opts, serverOptions...)

	srv := ssePlugin.NewServer(addr, opts...)

	// Post-construction setup callbacks (middleware, routes, streams).
	for _, setup := range serverSetups {
		setup(srv)
	}

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
