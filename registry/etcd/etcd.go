// Package etcd provides a bootstrap registry action for Etcd service registry.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/registry/etcd"
package etcd

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	etcdPlugin "github.com/tx7do/go-wind-plugins/registry/etcd"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	"github.com/tx7do/go-wind-bootstrap/tlsutil"
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
	if dialTimeout := c.GetDialTimeout(); dialTimeout > 0 {
		clientCfg.DialTimeout = time.Duration(dialTimeout) * time.Second
	}
	if syncInterval := c.GetAutoSyncInterval(); syncInterval > 0 {
		clientCfg.AutoSyncInterval = time.Duration(syncInterval) * time.Second
	}
	if keepAliveTime := c.GetDialKeepAliveTime(); keepAliveTime > 0 {
		clientCfg.DialKeepAliveTime = time.Duration(keepAliveTime) * time.Second
	}
	if keepAliveTimeout := c.GetDialKeepAliveTimeout(); keepAliveTimeout > 0 {
		clientCfg.DialKeepAliveTimeout = time.Duration(keepAliveTimeout) * time.Second
	}
	if c.GetRejectOldCluster() {
		clientCfg.RejectOldCluster = true
	}
	if c.GetPermitWithoutStream() {
		clientCfg.PermitWithoutStream = true
	}
	if tlsCfg, err := tlsutil.ClientTLS(c.GetTls()); err != nil {
		return nil, fmt.Errorf("etcd: load tls: %w", err)
	} else if tlsCfg != nil {
		clientCfg.TLS = tlsCfg
	}

	client, err := clientv3.New(clientCfg)
	if err != nil {
		return nil, fmt.Errorf("etcd: create client: %w", err)
	}

	var opts []etcdPlugin.Option
	if prefix := c.GetPrefix(); prefix != "" {
		opts = append(opts, etcdPlugin.Namespace(prefix))
	}
	if ttl := c.GetTtl(); ttl > 0 {
		opts = append(opts, etcdPlugin.RegisterTTL(time.Duration(ttl) * time.Second))
	}
	if maxRetry := c.GetMaxRetry(); maxRetry > 0 {
		opts = append(opts, etcdPlugin.MaxRetry(int(maxRetry)))
	}

	reg := etcdPlugin.New(client, opts...)

	regCleanup, err := bootstrap.RegisterInstance(ctx, reg, appCfg, endpoints)
	if err != nil {
		_ = client.Close()
		return nil, err
	}

	cleanup := func() {
		regCleanup()
		_ = client.Close()
	}
	return cleanup, nil
}
