// Package influxdb provides a bootstrap database builder for InfluxDB.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/database/influxdb"
package influxdb

import (
	"context"
	"fmt"

	influxdbCrud "github.com/tx7do/go-crud/influxdb"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterDatabaseBuilder(bootstrap.DatabaseTypeInfluxdb, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Database) (any, func(), error) {
	c := cfg.GetInfluxdb()
	if c == nil {
		return nil, nil, fmt.Errorf("influxdb: config is nil")
	}

	// 全部配置字段 → 选项的映射由 protoc-gen-fieldmap 生成（mapping.spec.yaml）。
	options := MapInfluxDBOptions(c, nil)

	client, err := influxdbCrud.NewClient(options...)
	if err != nil {
		return nil, nil, fmt.Errorf("influxdb: create client failed: %w", err)
	}

	return client, func() {}, nil
}
