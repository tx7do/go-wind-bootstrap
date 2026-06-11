// Package http provides a bootstrap server builder for the HTTP transport.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/transport/http"
//
// By default the builder creates a server with the std (net/http) driver.
// To configure the server (middleware + routes), call RegisterServerSetup before
// [bootstrap.Bootstrap]:
//
//	httpAdapter.RegisterServerSetup(func(srv *httpAdapter.Server) {
//	    srv.GET("/", myHandler)
//	})
package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	tokenbucket "github.com/tx7do/go-wind-plugins/ratelimit/tokenbucket"
	httpPlugin "github.com/tx7do/go-wind-plugins/transport/http"
	corsMW "github.com/tx7do/go-wind-plugins/transport/http/middleware/cors"
	loggingMW "github.com/tx7do/go-wind-plugins/transport/http/middleware/logging"
	ratelimitMW "github.com/tx7do/go-wind-plugins/transport/http/middleware/ratelimit"
	recoveryMW "github.com/tx7do/go-wind-plugins/transport/http/middleware/recovery"
	requestidMW "github.com/tx7do/go-wind-plugins/transport/http/middleware/requestid"
	timeoutMW "github.com/tx7do/go-wind-plugins/transport/http/middleware/timeout"
	tracingMW "github.com/tx7do/go-wind-plugins/transport/http/middleware/tracing"
	"github.com/tx7do/go-wind/transport"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

// Server is a type alias for the plugins HTTP Server.
// Code generators and user code should reference this type so that
// they only need to depend on the adapter package.
type Server = httpPlugin.Server

// serverSetups holds callbacks for configuring the HTTP server.
// They are called once during server construction, after the server is created
// but before it is returned to the bootstrap framework.
// Use them to register middleware (srv.Use) and routes (srv.GET, srv.POST, etc.).
var serverSetups []func(srv *Server)

// SetServerSetup sets the callback that will be invoked when the HTTP server
// is created. It replaces any previously registered callbacks.
//
// Deprecated: Use [RegisterServerSetup] for additive registration.
//
// This function must be called before [bootstrap.Bootstrap].
func SetServerSetup(fn func(srv *Server)) {
	serverSetups = []func(srv *Server){fn}
}

// RegisterServerSetup appends a callback that will be invoked when the HTTP
// server is created. All registered callbacks are called in registration order.
//
// This is the recommended API for code generators:
//
//	// generated_routes.go (auto-generated)
//	func init() {
//	    httpAdapter.RegisterServerSetup(registerProtoRoutes)
//	}
//
// This function must be called before [bootstrap.Bootstrap].
func RegisterServerSetup(fn func(srv *Server)) {
	serverSetups = append(serverSetups, fn)
}

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeHTTP, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	httpCfg := cfg.GetHttp()
	if httpCfg == nil {
		return nil, fmt.Errorf("http: config is nil")
	}

	addr := httpCfg.GetAddr()
	if addr == "" {
		addr = ":8080"
	}

	srv := httpPlugin.NewServer(addr, httpPlugin.WithDriver(newStdDriver()))

	// Register middleware from config.
	if mw := httpCfg.GetMiddleware(); mw != nil {
		applyMiddleware(srv, mw)
	}

	// Invoke all registered setup callbacks (middleware + routes).
	for _, setup := range serverSetups {
		setup(srv)
	}

	return srv, nil
}

// applyMiddleware reads the declarative middleware config and registers each
// enabled middleware on the server. The order follows the field number order
// in the proto definition.
func applyMiddleware(srv *httpPlugin.Server, mw *v1.Server_Http_Middleware) {
	if r := mw.GetRecovery(); r != nil {
		srv.Use(recoveryMW.Middleware(recoveryMW.WithStackTrace(r.GetStackTrace())))
	}
	if c := mw.GetCors(); c != nil {
		var opts []corsMW.Option
		if len(c.GetAllowedOrigins()) > 0 {
			opts = append(opts, corsMW.WithAllowedOrigins(c.GetAllowedOrigins()...))
		}
		if len(c.GetAllowedMethods()) > 0 {
			opts = append(opts, corsMW.WithAllowedMethods(c.GetAllowedMethods()...))
		}
		if len(c.GetAllowedHeaders()) > 0 {
			opts = append(opts, corsMW.WithAllowedHeaders(c.GetAllowedHeaders()...))
		}
		if len(c.GetExposedHeaders()) > 0 {
			opts = append(opts, corsMW.WithExposedHeaders(c.GetExposedHeaders()...))
		}
		if c.GetAllowCredentials() {
			opts = append(opts, corsMW.WithAllowCredentials(true))
		}
		if c.GetMaxAge() > 0 {
			opts = append(opts, corsMW.WithMaxAge(int(c.GetMaxAge())))
		}
		srv.Use(corsMW.Middleware(opts...))
	}
	if l := mw.GetLogging(); l != nil {
		var opts []loggingMW.Option
		if len(l.GetSkipPaths()) > 0 {
			opts = append(opts, loggingMW.WithSkipPaths(l.GetSkipPaths()...))
		}
		srv.Use(loggingMW.Middleware(opts...))
	}
	if rid := mw.GetRequestId(); rid != nil {
		headerName := rid.GetHeaderName()
		if headerName != "" {
			srv.Use(requestidMW.Middleware(requestidMW.WithHeaderName(headerName)))
		} else {
			srv.Use(requestidMW.Middleware())
		}
	}
	if t := mw.GetTracing(); t != nil {
		srv.Use(tracingMW.Middleware())
	}
	if rl := mw.GetRateLimit(); rl != nil {
		limiter, err := tokenbucket.New(float64(rl.GetRate()), float64(rl.GetBurst()))
		if err == nil {
			var rlOpts []ratelimitMW.Option
			if rl.GetWait() {
				rlOpts = append(rlOpts, ratelimitMW.WithWait())
			}
			srv.Use(ratelimitMW.Middleware(limiter, rlOpts...))
		}
	}
	if to := mw.GetTimeout(); to != nil {
		timeoutMs := to.GetDefaultTimeoutMs()
		if timeoutMs <= 0 {
			timeoutMs = 5000
		}
		srv.Use(timeoutMW.Middleware(time.Duration(timeoutMs) * time.Millisecond))
	}
}

// ---------------------------------------------------------------------------
// 内置 std Driver（基于 net/http 标准库）
// ---------------------------------------------------------------------------

type stdDriver struct {
	mux    *http.ServeMux
	server *http.Server
}

func newStdDriver() httpPlugin.Driver {
	return &stdDriver{mux: http.NewServeMux()}
}

func (d *stdDriver) Handle(method, path string, handler http.HandlerFunc) {
	d.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	})
}

func (d *stdDriver) Start(ctx context.Context, ln net.Listener) error {
	d.server = &http.Server{Handler: d.mux}
	errChan := make(chan error, 1)
	go func() {
		if err := d.server.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
			return
		}
		errChan <- nil
	}()
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return d.server.Shutdown(context.Background())
	}
}

func (d *stdDriver) Stop(ctx context.Context) error {
	if d.server == nil {
		return nil
	}
	return d.server.Shutdown(ctx)
}
