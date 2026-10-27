// Package elasticsearch provides a bootstrap database builder for Elasticsearch.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/database/elasticsearch"
package elasticsearch

import (
	"context"
	"fmt"

	esCrud "github.com/tx7do/go-crud/elasticsearch"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterDatabaseBuilder(bootstrap.DatabaseTypeElasticsearch, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Database) (any, func(), error) {
	c := cfg.GetElasticsearch()
	if c == nil {
		return nil, nil, fmt.Errorf("elasticsearch: config is nil")
	}

	// 全部配置字段 → 选项的映射由 protoc-gen-fieldmap 生成（mapping.spec.yaml）。
	options := MapElasticsearchOptions(c, nil)

	client, err := esCrud.NewElasticsearchClient(options...)
	if err != nil {
		return nil, nil, fmt.Errorf("elasticsearch: create client failed: %w", err)
	}

	return client, func() {}, nil
}
