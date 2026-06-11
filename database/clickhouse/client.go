// Package clickhouse provides a bootstrap database builder for ClickHouse.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/database/clickhouse"
package clickhouse

import (
	"context"
	"fmt"
	"time"

	clickhouseCrud "github.com/tx7do/go-crud/clickhouse"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterDatabaseBuilder(bootstrap.DatabaseTypeClickhouse, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Database) (any, func(), error) {
	c := cfg.GetClickhouse()
	if c == nil {
		return nil, nil, fmt.Errorf("clickhouse: config is nil")
	}

	var options []clickhouseCrud.Option

	if dsn := c.GetDsn(); dsn != "" {
		options = append(options, clickhouseCrud.WithDsn(dsn))
	}
	if addrs := c.GetAddresses(); len(addrs) > 0 {
		options = append(options, clickhouseCrud.WithAddresses(addrs...))
	}
	if db := c.GetDatabase(); db != "" {
		options = append(options, clickhouseCrud.WithDatabase(db))
	}
	if u := c.GetUsername(); u != "" {
		options = append(options, clickhouseCrud.WithUsername(u))
	}
	if p := c.GetPassword(); p != "" {
		options = append(options, clickhouseCrud.WithPassword(p))
	}
	options = append(options, clickhouseCrud.WithDebug(c.GetDebug()))

	if v := c.GetMaxOpenConns(); v > 0 {
		options = append(options, clickhouseCrud.WithMaxOpenConns(int(v)))
	}
	if v := c.GetMaxIdleConns(); v > 0 {
		options = append(options, clickhouseCrud.WithMaxIdleConns(int(v)))
	}
	if v := c.GetCompressionMethod(); v != "" {
		options = append(options, clickhouseCrud.WithCompressionMethod(v))
	}
	if v := c.GetDialTimeoutSeconds(); v > 0 {
		options = append(options, clickhouseCrud.WithDialTimeout(time.Duration(v)*time.Second))
	}
	if v := c.GetReadTimeoutSeconds(); v > 0 {
		options = append(options, clickhouseCrud.WithReadTimeout(time.Duration(v)*time.Second))
	}
	if v := c.GetConnMaxLifetimeSeconds(); v > 0 {
		options = append(options, clickhouseCrud.WithConnMaxLifetime(time.Duration(v)*time.Second))
	}
	if v := c.GetScheme(); v != "" {
		options = append(options, clickhouseCrud.WithScheme(v))
	}

	client, err := clickhouseCrud.NewClient(options...)
	if err != nil {
		return nil, nil, fmt.Errorf("clickhouse: create client failed: %w", err)
	}

	cleanup := func() { client.Close() }
	return client, cleanup, nil
}
