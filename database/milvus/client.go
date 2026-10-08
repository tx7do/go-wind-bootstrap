// Package milvus provides a bootstrap database builder for Milvus.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/database/milvus"
package milvus

import (
	"context"
	"fmt"

	milvusCrud "github.com/tx7do/go-crud/milvus"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterDatabaseBuilder(bootstrap.DatabaseTypeMilvus, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Database) (any, func(), error) {
	c := cfg.GetMilvus()
	if c == nil {
		return nil, nil, fmt.Errorf("milvus: config is nil")
	}

	var options []milvusCrud.Option

	if addr := c.GetAddress(); addr != "" {
		options = append(options, milvusCrud.WithAddress(addr))
	}
	if key := c.GetApiKey(); key != "" {
		// Zilliz Cloud：API Key 优先于用户名密码。
		options = append(options, milvusCrud.WithAPIKey(key))
	} else if u := c.GetUsername(); u != "" {
		options = append(options, milvusCrud.WithUsername(u))
		options = append(options, milvusCrud.WithPassword(c.GetPassword()))
	}
	if db := c.GetDbName(); db != "" {
		options = append(options, milvusCrud.WithDBName(db))
	}

	client, err := milvusCrud.NewClient(options...)
	if err != nil {
		return nil, nil, fmt.Errorf("milvus: create client failed: %w", err)
	}

	cleanup := func() { _ = client.Close() }
	return client, cleanup, nil
}
