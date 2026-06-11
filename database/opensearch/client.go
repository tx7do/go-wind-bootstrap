// Package opensearch provides a bootstrap database builder for OpenSearch.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/database/opensearch"
package opensearch

import (
	"context"
	"fmt"
	"time"

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

	var options []osCrud.Option

	if addrs := c.GetAddresses(); len(addrs) > 0 {
		options = append(options, osCrud.WithAddresses(addrs...))
	}
	if u := c.GetUsername(); u != "" {
		options = append(options, osCrud.WithUsername(u))
	}
	if p := c.GetPassword(); p != "" {
		options = append(options, osCrud.WithPassword(p))
	}

	options = append(options, osCrud.WithEnableMetrics(c.GetEnableMetrics()))
	options = append(options, osCrud.WithEnableDebugLogger(c.GetEnableDebugLogger()))
	options = append(options, osCrud.WithCompressRequestBody(c.GetCompressRequestBody()))
	options = append(options, osCrud.WithDiscoverNodesOnStart(c.GetDiscoverNodesOnStart()))

	if v := c.GetMaxRetries(); v > 0 {
		options = append(options, osCrud.WithMaxRetries(int(v)))
	}
	if v := c.GetDiscoverNodesIntervalSeconds(); v > 0 {
		options = append(options, osCrud.WithDiscoverNodesInterval(time.Duration(v)*time.Second))
	}

	client, err := osCrud.NewOpenSearchClient(options...)
	if err != nil {
		return nil, nil, fmt.Errorf("opensearch: create client failed: %w", err)
	}

	return client, func() {}, nil
}
