// Package opensearch provides a bootstrap database builder for OpenSearch.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/database/opensearch"
package opensearch

import (
	"context"
	"fmt"

	osCrud "github.com/tx7do/go-crud/opensearch"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterDatabaseBuilder(bootstrap.DatabaseTypeOpensearch, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Database) (any, func(), error) {
	c := cfg.GetOpensearch()
	if c == nil {
		return nil, nil, fmt.Errorf("opensearch: config is nil")
	}

	// 全部配置字段 → 选项的映射由 protoc-gen-fieldmap 生成（mapping.spec.yaml）。
	options := MapOpenSearchOptions(c, nil)

	client, err := osCrud.NewOpenSearchClient(options...)
	if err != nil {
		return nil, nil, fmt.Errorf("opensearch: create client failed: %w", err)
	}

	return client, func() {}, nil
}
