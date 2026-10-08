// Package neo4j provides a bootstrap database builder for Neo4j.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/database/neo4j"
package neo4j

import (
	"context"
	"fmt"

	neo4jCrud "github.com/tx7do/go-crud/neo4j"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterDatabaseBuilder(bootstrap.DatabaseTypeNeo4j, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Database) (any, func(), error) {
	c := cfg.GetNeo4J()
	if c == nil {
		return nil, nil, fmt.Errorf("neo4j: config is nil")
	}

	var options []neo4jCrud.Option

	if uri := c.GetUri(); uri != "" {
		options = append(options, neo4jCrud.WithURI(uri))
	}
	if u := c.GetUsername(); u != "" {
		options = append(options, neo4jCrud.WithBasicAuth(u, c.GetPassword(), c.GetRealm()))
	}

	client, err := neo4jCrud.NewClient(options...)
	if err != nil {
		return nil, nil, fmt.Errorf("neo4j: create client failed: %w", err)
	}

	cleanup := func() { _ = client.Close() }
	return client, cleanup, nil
}
