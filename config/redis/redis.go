// Package redis provides a bootstrap config action for Redis config source.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/config/redis"
package redis

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"

	redisPlugin "github.com/tx7do/go-wind-plugins/config/redis"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterConfigAction(bootstrap.ConfigTypeRedis, newAction)
}

func newAction(ctx context.Context, cfg *v1.Config) (func(), error) {
	c := cfg.GetRedis()
	if c == nil {
		return nil, fmt.Errorf("redis: config is nil")
	}

	// 连接字段 → 客户端配置的映射由 protoc-gen-fieldmap 生成（mapping.spec.yaml）。
	rdb := goredis.NewClient(NewRedisConfig(c))

	// path 为插件选项，保留手工段。
	var opts []redisPlugin.Option
	if path := c.GetPath(); path != "" {
		opts = append(opts, redisPlugin.WithPath(path))
	}

	_, err := redisPlugin.New(rdb, opts...)
	if err != nil {
		rdb.Close()
		return nil, fmt.Errorf("redis: create source: %w", err)
	}

	cleanup := func() {
		rdb.Close()
	}
	return cleanup, nil
}
