// Package doris provides a bootstrap database builder for Apache Doris.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/database/doris"
package doris

import (
	"context"
	"fmt"

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

	// 常规配置字段 → 选项的映射由 protoc-gen-fieldmap 生成（mapping.spec.yaml）。
	var options []dorisCrud.Option
	options = MapDorisOptions(c, options)

	// Stream Load 认证合成一项，保留手工段。
	if u := c.GetStreamLoadUsername(); u != "" && c.GetStreamLoadPassword() != "" {
		options = append(options, dorisCrud.WithStreamLoadAuth(u, c.GetStreamLoadPassword()))
	}

	client, err := dorisCrud.NewClient(options...)
	if err != nil {
		return nil, nil, fmt.Errorf("doris: create client failed: %w", err)
	}

	cleanup := func() { _ = client.Close() }
	return client, cleanup, nil
}
