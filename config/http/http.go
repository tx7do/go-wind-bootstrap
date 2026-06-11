// Package http provides a bootstrap config action for HTTP-based config source.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/config/http"
package http

import (
	"context"
	"fmt"

	httpPlugin "github.com/tx7do/go-wind-plugins/config/http"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterConfigAction(bootstrap.ConfigTypeHTTP, newAction)
}

func newAction(ctx context.Context, cfg *v1.Config) (func(), error) {
	c := cfg.GetHttp()
	if c == nil {
		return nil, fmt.Errorf("http: config is nil")
	}

	var opts []httpPlugin.Option

	if u := c.GetUrl(); u != "" {
		opts = append(opts, httpPlugin.WithURL(u))
	}
	if m := c.GetMethod(); m != "" {
		opts = append(opts, httpPlugin.WithMethod(m))
	}
	for k, v := range c.GetHeaders() {
		opts = append(opts, httpPlugin.WithHeader(k, v))
	}

	src, err := httpPlugin.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("http: create source: %w", err)
	}

	cleanup := func() {
		_ = src.Close()
	}
	return cleanup, nil
}
