// Package cassandra provides a bootstrap database builder for Cassandra.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/database/cassandra"
package cassandra

import (
	"context"
	"fmt"
	"time"

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

	var options []cassandraCrud.Option

	if addr := c.GetAddress(); addr != "" {
		options = append(options, cassandraCrud.WithHosts(addr))
	}
	if u := c.GetUsername(); u != "" {
		options = append(options, cassandraCrud.WithUsername(u))
		options = append(options, cassandraCrud.WithPassword(c.GetPassword()))
	}
	if ks := c.GetKeyspace(); ks != "" {
		options = append(options, cassandraCrud.WithKeyspace(ks))
	}
	// 一致性级别：未配置时由 go-crud 缺省为 Quorum
	// （gocql 零值 Any 仅对写合法，服务端会拒绝全部读路径）。
	if v := c.GetConsistency(); v > 0 {
		options = append(options, cassandraCrud.WithConsistency(uint32(v)))
	}
	if v := c.GetConnectTimeoutSeconds(); v > 0 {
		options = append(options, cassandraCrud.WithConnectTimeout(time.Duration(v)*time.Second))
	}
	if v := c.GetTimeoutSeconds(); v > 0 {
		options = append(options, cassandraCrud.WithTimeout(time.Duration(v)*time.Second))
	}
	options = append(options, cassandraCrud.WithDisableInitialHostLookup(c.GetDisableInitialHostLookup()))

	client, err := cassandraCrud.NewCassandraClient(options...)
	if err != nil {
		return nil, nil, fmt.Errorf("cassandra: create client failed: %w", err)
	}

	cleanup := func() { client.Close() }
	return client, cleanup, nil
}
