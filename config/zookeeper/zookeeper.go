// Package zookeeper provides a bootstrap config action for Zookeeper config source.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/config/zookeeper"
package zookeeper

import (
	"context"
	"fmt"
	"time"

	"github.com/go-zookeeper/zk"

	zkPlugin "github.com/tx7do/go-wind-plugins/config/zookeeper"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterConfigAction(bootstrap.ConfigTypeZookeeper, newAction)
}

func newAction(ctx context.Context, cfg *v1.Config) (func(), error) {
	c := cfg.GetZookeeper()
	if c == nil {
		return nil, fmt.Errorf("zookeeper: config is nil")
	}

	endpoints := c.GetEndpoints()
	if len(endpoints) == 0 {
		return nil, fmt.Errorf("zookeeper: no endpoints")
	}

	conn, _, err := zk.Connect(endpoints, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("zookeeper: connect: %w", err)
	}

	var opts []zkPlugin.Option
	if path := c.GetPath(); path != "" {
		opts = append(opts, zkPlugin.WithPath(path))
	}

	_, err = zkPlugin.New(conn, opts...)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("zookeeper: create source: %w", err)
	}

	cleanup := func() {
		conn.Close()
	}
	return cleanup, nil
}
