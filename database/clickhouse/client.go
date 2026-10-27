// Package clickhouse provides a bootstrap database builder for ClickHouse.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/database/clickhouse"
package clickhouse

import (
	"context"
	"fmt"

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

	// 全部配置字段 → 选项的映射由 protoc-gen-fieldmap 生成（mapping.spec.yaml）。
	options := MapClickHouseOptions(c, nil)

	client, err := clickhouseCrud.NewClient(options...)
	if err != nil {
		return nil, nil, fmt.Errorf("clickhouse: create client failed: %w", err)
	}

	cleanup := func() { client.Close() }
	return client, cleanup, nil
}
