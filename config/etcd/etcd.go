// Package etcd provides a bootstrap config action for Etcd config source.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/config/etcd"
package etcd

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"

	etcdPlugin "github.com/tx7do/go-wind-plugins/config/etcd"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterConfigAction(bootstrap.ConfigTypeEtcd, newAction)
}

func newAction(ctx context.Context, cfg *v1.Config) (func(), error) {
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

	var opts []etcdPlugin.Option
	if path := c.GetPath(); path != "" {
		opts = append(opts, etcdPlugin.WithPath(path))
	}
	if c.GetPrefix() {
		opts = append(opts, etcdPlugin.WithPrefix(true))
	}

	_, err = etcdPlugin.New(client, opts...)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("etcd: create source: %w", err)
	}

	cleanup := func() {
		client.Close()
	}
	return cleanup, nil
}
