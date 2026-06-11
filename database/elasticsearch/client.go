// Package elasticsearch provides a bootstrap database builder for Elasticsearch.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/database/elasticsearch"
package elasticsearch

import (
	"context"
	"fmt"
	"time"

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

	var options []esCrud.Option

	if addrs := c.GetAddresses(); len(addrs) > 0 {
		options = append(options, esCrud.WithAddresses(addrs...))
	}
	if u := c.GetUsername(); u != "" {
		options = append(options, esCrud.WithUsername(u))
	}
	if p := c.GetPassword(); p != "" {
		options = append(options, esCrud.WithPassword(p))
	}
	if id := c.GetCloudId(); id != "" {
		options = append(options, esCrud.WithCloudID(id))
	}
	if key := c.GetApiKey(); key != "" {
		options = append(options, esCrud.WithAPIKey(key))
	}
	if st := c.GetServiceToken(); st != "" {
		options = append(options, esCrud.WithServiceToken(st))
	}
	if fp := c.GetCertificateFingerprint(); fp != "" {
		options = append(options, esCrud.WithCertificateFingerprint(fp))
	}

	options = append(options, esCrud.WithEnableMetrics(c.GetEnableMetrics()))
	options = append(options, esCrud.WithEnableDebugLogger(c.GetEnableDebugLogger()))
	options = append(options, esCrud.WithCompressRequestBody(c.GetCompressRequestBody()))
	options = append(options, esCrud.WithDiscoverNodesOnStart(c.GetDiscoverNodesOnStart()))

	if v := c.GetMaxRetries(); v > 0 {
		options = append(options, esCrud.WithMaxRetries(int(v)))
	}
	if v := c.GetDiscoverNodesIntervalSeconds(); v > 0 {
		options = append(options, esCrud.WithDiscoverNodesInterval(time.Duration(v)*time.Second))
	}

	client, err := esCrud.NewElasticsearchClient(options...)
	if err != nil {
		return nil, nil, fmt.Errorf("elasticsearch: create client failed: %w", err)
	}

	return client, func() {}, nil
}
