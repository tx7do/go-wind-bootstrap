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

	"github.com/gocql/gocql"

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

	cluster := gocql.NewCluster(c.GetAddress())

	// 认证
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: c.GetUsername(),
		Password: c.GetPassword(),
	}

	// Keyspace
	if ks := c.GetKeyspace(); ks != "" {
		cluster.Keyspace = ks
	}

	// 一致性级别
	if v := c.GetConsistency(); v > 0 {
		cluster.Consistency = gocql.Consistency(v)
	}

	// 超时
	if v := c.GetConnectTimeoutSeconds(); v > 0 {
		cluster.ConnectTimeout = time.Duration(v) * time.Second
	}
	if v := c.GetTimeoutSeconds(); v > 0 {
		cluster.Timeout = time.Duration(v) * time.Second
	}

	// 禁止主机查找
	cluster.DisableInitialHostLookup = c.GetDisableInitialHostLookup()

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, nil, fmt.Errorf("cassandra: create session failed: %w", err)
	}

	cleanup := func() { session.Close() }
	return session, cleanup, nil
}
