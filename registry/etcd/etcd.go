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

func newAction(ctx context.Context, cfg *v1.Registry) (bootstrap.Registry, func(), error) {
	c := cfg.GetEtcd()
	if c == nil {
		return nil, nil, fmt.Errorf("etcd: config is nil")
	}

	// 常规字段（endpoints/username/…/permit_without_stream）由 FillEtcdConfig 生成映射。
	clientCfg := FillEtcdConfig(c, &clientv3.Config{})
	// TLS 装载经 tlsutil 转换，保留手工段。
	if tlsCfg, err := tlsutil.ClientTLS(c.GetTls()); err != nil {
		return nil, nil, fmt.Errorf("etcd: load tls: %w", err)
	} else if tlsCfg != nil {
		clientCfg.TLS = tlsCfg
	}

	client, err := clientv3.New(*clientCfg)
	if err != nil {
		return nil, nil, fmt.Errorf("etcd: create client: %w", err)
	}

	var opts []etcdPlugin.Option
	if prefix := c.GetPrefix(); prefix != "" {
		opts = append(opts, etcdPlugin.Namespace(prefix))
	}
	if ttl := c.GetTtl(); ttl > 0 {
		opts = append(opts, etcdPlugin.RegisterTTL(time.Duration(ttl)*time.Second))
	}
	if maxRetry := c.GetMaxRetry(); maxRetry > 0 {
		opts = append(opts, etcdPlugin.MaxRetry(int(maxRetry)))
	}

	reg := etcdPlugin.New(client, opts...)

	return reg, func() { _ = client.Close() }, nil
}
