// Package qdrant provides a bootstrap database builder for Qdrant.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/database/qdrant"
package qdrant

import (
	"context"
	"fmt"

	qdrantCrud "github.com/tx7do/go-crud/qdrant"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	"github.com/tx7do/go-wind-bootstrap/tlsutil"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterDatabaseBuilder(bootstrap.DatabaseTypeQdrant, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Database) (any, func(), error) {
	c := cfg.GetQdrant()
	if c == nil {
		return nil, nil, fmt.Errorf("qdrant: config is nil")
	}

	var options []qdrantCrud.Option

	if host := c.GetHost(); host != "" {
		options = append(options, qdrantCrud.WithHost(host))
	}
	if port := c.GetPort(); port > 0 {
		options = append(options, qdrantCrud.WithPort(int(port)))
	}
	if key := c.GetApiKey(); key != "" {
		options = append(options, qdrantCrud.WithAPIKey(key))
	}

	// TLS：显式开关优先；开启但未提供证书材料时交给官方 SDK 默认（系统根证书）。
	useTLS := c.GetUseTls()
	tlsCfg, err := tlsutil.ClientTLS(c.GetTls())
	if err != nil {
		return nil, nil, fmt.Errorf("qdrant: load tls: %w", err)
	}
	if tlsCfg != nil {
		options = append(options, qdrantCrud.WithTLSConfig(tlsCfg))
		useTLS = true
	}
	if useTLS {
		options = append(options, qdrantCrud.WithTLS(true))
	}

	client, err := qdrantCrud.NewClient(options...)
	if err != nil {
		return nil, nil, fmt.Errorf("qdrant: create client failed: %w", err)
	}

	cleanup := func() { _ = client.Close() }
	return client, cleanup, nil
}
