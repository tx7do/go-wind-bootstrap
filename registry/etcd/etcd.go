// Package etcd provides a bootstrap registry action for Etcd service registry.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/registry/etcd"
package etcd

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"

	etcdPlugin "github.com/tx7do/go-wind-plugins/registry/etcd"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterRegistryAction(bootstrap.RegistryTypeEtcd, newAction)
}

func newAction(ctx context.Context, appCfg *v1.App, endpoints []string, cfg *v1.Registry) (func(), error) {
	c := cfg.GetEtcd()
	if c == nil {
		return nil, fmt.Errorf("etcd: config is nil")
	}

	clientCfg := clientv3.Config{
		Endpoints: c.GetEndpoints(),
	}
	if username := c.GetUsername(); username != "" {
		clientCfg.Username = username
	}
	if password := c.GetPassword(); password != "" {
		clientCfg.Password = password
	}

	client, err := clientv3.New(clientCfg)
	if err != nil {
		return nil, fmt.Errorf("etcd: create client: %w", err)
	}

	reg := etcdPlugin.New(client)
	_ = reg

	cleanup := func() {
		client.Close()
	}
	return cleanup, nil
}
