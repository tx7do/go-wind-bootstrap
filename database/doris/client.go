// Package doris provides a bootstrap database builder for Apache Doris.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/database/doris"
package doris

import (
	"context"
	"fmt"
	"time"

	dorisCrud "github.com/tx7do/go-crud/doris"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterDatabaseBuilder(bootstrap.DatabaseTypeDoris, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Database) (any, func(), error) {
	c := cfg.GetDoris()
	if c == nil {
		return nil, nil, fmt.Errorf("doris: config is nil")
	}

	var options []dorisCrud.Option

	if dsn := c.GetDsn(); dsn != "" {
		options = append(options, dorisCrud.WithDSN(dsn))
	}
	if v := c.GetMaxIdleConnections(); v > 0 {
		options = append(options, dorisCrud.WithMaxIdleConns(int(v)))
	}
	if v := c.GetMaxOpenConnections(); v > 0 {
		options = append(options, dorisCrud.WithMaxOpenConns(int(v)))
	}
	if v := c.GetConnectionMaxLifetimeSeconds(); v > 0 {
		options = append(options, dorisCrud.WithConnMaxLifetime(time.Duration(v)*time.Second))
	}

	// Stream Load 配置
	if ep := c.GetStreamLoadEndpoint(); ep != "" {
		options = append(options, dorisCrud.WithStreamLoadEndpoint(ep))
	}
	if u := c.GetStreamLoadUsername(); u != "" && c.GetStreamLoadPassword() != "" {
		options = append(options, dorisCrud.WithStreamLoadAuth(u, c.GetStreamLoadPassword()))
	}
	if v := c.GetStreamLoadTimeoutSeconds(); v > 0 {
		options = append(options, dorisCrud.WithStreamLoadTimeout(time.Duration(v)*time.Second))
	}
	if m := c.GetStreamLoadMethod(); m != "" {
		options = append(options, dorisCrud.WithStreamLoadMethod(m))
	}

	client, err := dorisCrud.NewClient(options...)
	if err != nil {
		return nil, nil, fmt.Errorf("doris: create client failed: %w", err)
	}

	cleanup := func() { _ = client.Close() }
	return client, cleanup, nil
}
