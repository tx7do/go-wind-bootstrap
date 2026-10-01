// Package zookeeper provides a bootstrap registry action for Zookeeper service registry.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/registry/zookeeper"
package zookeeper

import (
	"context"
	"fmt"
	"time"

	"github.com/go-zookeeper/zk"

	zkPlugin "github.com/tx7do/go-wind-plugins/registry/zookeeper"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterRegistryAction(bootstrap.RegistryTypeZookeeper, newAction)
}

func newAction(ctx context.Context, appCfg *v1.App, endpoints []string, cfg *v1.Registry) (func(), error) {
	c := cfg.GetZookeeper()
	if c == nil {
		return nil, fmt.Errorf("zookeeper: config is nil")
	}

	eps := c.GetEndpoints()
	if len(eps) == 0 {
		return nil, fmt.Errorf("zookeeper: no endpoints")
	}

	sessionTimeout := 5 * time.Second
	if t := c.GetSessionTimeout(); t > 0 {
		sessionTimeout = time.Duration(t) * time.Second
	}
	conn, _, err := zk.Connect(eps, sessionTimeout)
	if err != nil {
		return nil, fmt.Errorf("zookeeper: connect: %w", err)
	}

	var opts []zkPlugin.Option
	if rootPath := c.GetRootPath(); rootPath != "" {
		opts = append(opts, zkPlugin.WithRootPath(rootPath))
	}
	if username := c.GetUsername(); username != "" {
		opts = append(opts, zkPlugin.WithDigestACL(username, c.GetPassword()))
	}

	reg := zkPlugin.New(conn, opts...)

	regCleanup, err := bootstrap.RegisterInstance(ctx, reg, appCfg, endpoints)
	if err != nil {
		conn.Close()
		return nil, err
	}

	cleanup := func() {
		regCleanup()
		conn.Close()
	}
	return cleanup, nil
}
