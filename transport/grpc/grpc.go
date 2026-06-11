// Package grpc provides a bootstrap server builder for the gRPC transport.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/transport/grpc"
//
// To register gRPC services and/or middleware, call the setup functions
// before [bootstrap.Bootstrap]:
//
//	grpcAdapter.SetMiddlewares(myInterceptor)
//	grpcAdapter.SetServiceRegistrar(func(srv *grpc.Server) {
//	    pb.RegisterGreeterServer(srv, &greeterService{})
//	})
package grpc

import (
	"fmt"

	grpcLogging "github.com/tx7do/go-wind-plugins/transport/grpc/middleware/logging"
	grpcRecovery "github.com/tx7do/go-wind-plugins/transport/grpc/middleware/recovery"
	grpcTracing "github.com/tx7do/go-wind-plugins/transport/grpc/middleware/tracing"
	grpcValidate "github.com/tx7do/go-wind-plugins/transport/grpc/middleware/validate"
	grpcPlugin "github.com/tx7do/go-wind-plugins/transport/grpc/server"
	"github.com/tx7do/go-wind/transport"
	"google.golang.org/grpc"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

// unaryMiddlewares holds user-provided gRPC unary interceptors.
var unaryMiddlewares []grpc.UnaryServerInterceptor

// SetMiddlewares appends gRPC unary interceptors that will be applied to the
// server. Call before [bootstrap.Bootstrap].
func SetMiddlewares(middlewares ...grpc.UnaryServerInterceptor) {
	unaryMiddlewares = append(unaryMiddlewares, middlewares...)
}

// serviceRegistrars holds callbacks for registering gRPC services.
var serviceRegistrars []func(srv *grpc.Server)

// SetServiceRegistrar sets the callback for registering gRPC services.
// It replaces any previously registered callbacks.
//
// Deprecated: Use [RegisterServiceRegistrar] for additive registration.
//
// This function must be called before [bootstrap.Bootstrap].
func SetServiceRegistrar(fn func(srv *grpc.Server)) {
	serviceRegistrars = []func(srv *grpc.Server){fn}
}

// RegisterServiceRegistrar appends a callback for registering gRPC services.
// All registered callbacks are called in registration order.
//
// This is the recommended API for code generators:
//
//	// generated_grpc.go (auto-generated)
//	func init() {
//	    grpcAdapter.RegisterServiceRegistrar(registerProtoServices)
//	}
//
// This function must be called before [bootstrap.Bootstrap].
func RegisterServiceRegistrar(fn func(srv *grpc.Server)) {
	serviceRegistrars = append(serviceRegistrars, fn)
}

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeGRPC, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	grpcCfg := cfg.GetGrpc()
	if grpcCfg == nil {
		return nil, fmt.Errorf("grpc: config is nil")
	}

	addr := grpcCfg.GetAddr()
	if addr == "" {
		addr = ":9000"
	}

	var opts []grpcPlugin.Option

	// Build interceptor chain from config + programmatic middleware.
	var interceptors []grpc.UnaryServerInterceptor
	if mw := grpcCfg.GetMiddleware(); mw != nil {
		if mw.GetRecovery() != nil {
			interceptors = append(interceptors, grpcRecovery.UnaryInterceptor())
		}
		if mw.GetLogging() != nil {
			interceptors = append(interceptors, grpcLogging.UnaryInterceptor())
		}
		if mw.GetTracing() != nil {
			interceptors = append(interceptors, grpcTracing.UnaryInterceptor())
		}
		if mw.GetValidate() != nil {
			interceptors = append(interceptors, grpcValidate.UnaryServerInterceptor())
		}
	}
	interceptors = append(interceptors, unaryMiddlewares...)

	if len(interceptors) > 0 {
		opts = append(opts, grpcPlugin.WithMiddleware(interceptors...))
	}

	// Register user services on the underlying grpc.Server.
	if len(serviceRegistrars) > 0 {
		rawSrv := grpc.NewServer()
		for _, registrar := range serviceRegistrars {
			registrar(rawSrv)
		}
		opts = append(opts, grpcPlugin.WithServer(rawSrv))
	}

	srv := grpcPlugin.NewServer(addr, opts...)
	return srv, nil
}
