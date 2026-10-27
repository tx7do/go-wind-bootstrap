// Package cassandra provides a bootstrap database builder for Cassandra.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/database/cassandra"
package cassandra

import (
	"context"
	"fmt"

	cassandraCrud "github.com/tx7do/go-crud/cassandra"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterDatabaseBuilder(bootstrap.DatabaseTypeCassandra, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Database) (any, func(), error) {
	c := cfg.GetCassandra()
	if c == nil {
		return nil, nil, fmt.Errorf("cassandra: config is nil")
	}

	// 常规配置字段 → 选项的映射由 protoc-gen-fieldmap 生成（mapping.spec.yaml）。
	var options []cassandraCrud.Option
	options = MapCassandraOptions(c, options)

	// 用户名密码成对写入，保留手工段。
	if u := c.GetUsername(); u != "" {
		options = append(options, cassandraCrud.WithUsername(u))
		options = append(options, cassandraCrud.WithPassword(c.GetPassword()))
	}

	client, err := cassandraCrud.NewCassandraClient(options...)
	if err != nil {
		return nil, nil, fmt.Errorf("cassandra: create client failed: %w", err)
	}

	cleanup := func() { client.Close() }
	return client, cleanup, nil
}
