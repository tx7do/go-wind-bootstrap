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
	"github.com/tx7do/go-wind-bootstrap/tlsutil"
)

func init() {
	bootstrap.MustRegisterConfigAction(bootstrap.ConfigTypeEtcd, newAction)
}

func newAction(ctx context.Context, cfg *v1.Config) (func(), error) {
	c := cfg.GetEtcd()
	if c == nil {
		return nil, fmt.Errorf("etcd: config is nil")
	}

	// 常规连接字段 → 客户端配置的映射由 protoc-gen-fieldmap 生成（mapping.spec.yaml）。
	clientCfg := clientv3.Config{}
	FillEtcdConfig(c, &clientCfg)

	// TLS 证书装载经 tlsutil 转换，保留手工段。
	if tlsCfg, err := tlsutil.ClientTLS(c.GetTls()); err != nil {
		return nil, fmt.Errorf("etcd: load tls: %w", err)
	} else if tlsCfg != nil {
		clientCfg.TLS = tlsCfg
	}

	client, err := clientv3.New(clientCfg)
	if err != nil {
		return nil, fmt.Errorf("etcd: create client: %w", err)
	}

	// path/prefix 为插件选项，保留手工段。
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
