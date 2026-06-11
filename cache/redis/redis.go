// Package redis provides a bootstrap cache builder for Redis-backed caching.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/cache/redis"
package redis

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"

	redisPlugin "github.com/tx7do/go-wind-plugins/cache/redis"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterCacheBuilder(bootstrap.CacheTypeRedis, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Cache) (any, func(), error) {
	c := cfg.GetRedis()
	if c == nil {
		return nil, nil, fmt.Errorf("redis cache: config is nil")
	}

	addr := c.GetAddr()
	if addr == "" {
		addr = "localhost:6379"
	}

	rdb := goredis.NewClient(&goredis.Options{
		Addr:     addr,
		Password: c.GetPassword(),
		DB:       int(c.GetDb()),
	})

	// Verify connectivity.
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, nil, fmt.Errorf("redis cache: ping failed: %w", err)
	}

	var opts []redisPlugin.Option
	if prefix := c.GetKeyPrefix(); prefix != "" {
		opts = append(opts, redisPlugin.WithKeyPrefix(prefix))
	}

	cache := redisPlugin.New(rdb, opts...)

	cleanup := func() {
		_ = cache.Close()
		_ = rdb.Close()
	}

	return cache, cleanup, nil
}
