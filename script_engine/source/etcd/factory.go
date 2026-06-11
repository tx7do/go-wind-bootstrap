// Package etcd provides etcd KV as a script source.
package etcd

import (
	"context"
	"fmt"

	"github.com/tx7do/go-scripts/source"
	etcdSource "github.com/tx7do/go-scripts/source/etcd"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	scriptEngine "github.com/tx7do/go-wind-bootstrap/script_engine"
)

func init() {
	scriptEngine.MustRegisterSourceFactory(v1.Script_Source_ETCD, NewSource)
}

// NewSource 根据配置创建 etcd 脚本来源。
func NewSource(cfg *v1.Script_Source) (source.Reader, error) {
	if cfg == nil {
		return nil, fmt.Errorf("etcd source: config is nil")
	}

	opts := cfg.GetEtcdOptions()
	if opts == nil {
		return nil, fmt.Errorf("etcd source: etcd_options is required")
	}

	var etcdOpts []etcdSource.Option

	if len(opts.GetEndpoints()) > 0 {
		etcdOpts = append(etcdOpts, etcdSource.WithEndpoints(opts.GetEndpoints()...))
	}
	if prefix := opts.GetPrefix(); prefix != "" {
		etcdOpts = append(etcdOpts, etcdSource.WithPrefix(prefix))
	}

	return etcdSource.New(context.Background(), etcdOpts...)
}
