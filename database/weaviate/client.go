// Package weaviate provides a bootstrap database builder for Weaviate.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/database/weaviate"
package weaviate

import (
	"context"
	"fmt"

	weaviateCrud "github.com/tx7do/go-crud/weaviate"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterDatabaseBuilder(bootstrap.DatabaseTypeWeaviate, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Database) (any, func(), error) {
	c := cfg.GetWeaviate()
	if c == nil {
		return nil, nil, fmt.Errorf("weaviate: config is nil")
	}

	var options []weaviateCrud.Option

	if host := c.GetHost(); host != "" {
		options = append(options, weaviateCrud.WithHost(host))
	}
	if scheme := c.GetScheme(); scheme != "" {
		options = append(options, weaviateCrud.WithScheme(scheme))
	}

	// 认证：Weaviate 官方客户端以请求头承载（Authorization: Bearer <key>）。
	headers := map[string]string{}
	for k, v := range c.GetHeaders() {
		headers[k] = v
	}
	if key := c.GetApiKey(); key != "" {
		headers["Authorization"] = "Bearer " + key
	}
	if len(headers) > 0 {
		options = append(options, weaviateCrud.WithHeaders(headers))
	}

	client, err := weaviateCrud.NewClient(options...)
	if err != nil {
		return nil, nil, fmt.Errorf("weaviate: create client failed: %w", err)
	}

	cleanup := func() { _ = client.Close() }
	return client, cleanup, nil
}
