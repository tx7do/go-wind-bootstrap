// Package redis provides a bootstrap cache builder for Redis-backed caching.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/cache/redis"
package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"

	redisPlugin "github.com/tx7do/go-wind-plugins/cache/redis"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	"github.com/tx7do/go-wind-bootstrap/tlsutil"
)

func init() {
	bootstrap.MustRegisterCacheBuilder(bootstrap.CacheTypeRedis, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Cache) (any, func(), error) {
	c := cfg.GetRedis()
	if c == nil {
		return nil, nil, fmt.Errorf("redis cache: config is nil")
	}

	rdb, err := newClient(ctx, c)
	if err != nil {
		return nil, nil, err
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

// newClient builds a [goredis.UniversalClient] from the conf: standalone by
// default, cluster or sentinel when mode/addrs/master_name say so.
func newClient(ctx context.Context, c *v1.Cache_Redis) (goredis.UniversalClient, error) {
	tlsCfg, err := tlsutil.ClientTLS(c.GetTls())
	if err != nil {
		return nil, fmt.Errorf("redis cache: load tls: %w", err)
	}

	base := goredis.UniversalOptions{
		Username:        c.GetUsername(),
		Password:        c.GetPassword(),
		TLSConfig:       tlsCfg,
		DialTimeout:     seconds(c.GetDialTimeout()),
		ReadTimeout:     seconds(c.GetReadTimeout()),
		WriteTimeout:    seconds(c.GetWriteTimeout()),
		PoolTimeout:     seconds(c.GetPoolTimeout()),
		ConnMaxIdleTime: seconds(c.GetConnMaxIdleTime()),
		ConnMaxLifetime: seconds(c.GetConnMaxLifetime()),
		MinRetryBackoff: millis(c.GetMinRetryBackoffMs()),
		MaxRetryBackoff: millis(c.GetMaxRetryBackoffMs()),
		PoolSize:        int(c.GetPoolSize()),
		MinIdleConns:    int(c.GetMinIdleConns()),
		MaxRetries:      int(c.GetMaxRetries()),
	}

	var rdb goredis.UniversalClient

	switch mode := c.GetMode(); mode {
	case "", "standalone":
		addr := c.GetAddr()
		if addr == "" {
			addr = "localhost:6379"
		}
		base.Addrs = []string{addr}
		base.DB = int(c.GetDb())
		rdb = goredis.NewUniversalClient(&base)

	case "cluster":
		addrs := c.GetAddrs()
		if len(addrs) == 0 && c.GetAddr() != "" {
			addrs = []string{c.GetAddr()}
		}
		if len(addrs) == 0 {
			return nil, fmt.Errorf("redis cache: cluster mode requires addrs")
		}
		base.Addrs = addrs
		rdb = goredis.NewUniversalClient(&base)

	case "sentinel":
		if c.GetMasterName() == "" {
			return nil, fmt.Errorf("redis cache: sentinel mode requires master_name")
		}
		addrs := c.GetAddrs()
		if len(addrs) == 0 && c.GetAddr() != "" {
			addrs = []string{c.GetAddr()}
		}
		if len(addrs) == 0 {
			return nil, fmt.Errorf("redis cache: sentinel mode requires addrs")
		}
		base.Addrs = addrs
		base.MasterName = c.GetMasterName()
		base.SentinelUsername = c.GetSentinelUsername()
		base.SentinelPassword = c.GetSentinelPassword()
		rdb = goredis.NewUniversalClient(&base)

	default:
		return nil, fmt.Errorf("redis cache: unknown mode %q (want standalone/cluster/sentinel)", mode)
	}

	// Verify connectivity (cluster clients tolerate partial nodes failing).
	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("redis cache: ping failed: %w", err)
	}
	return rdb, nil
}

func seconds(v int32) time.Duration {
	return time.Duration(v) * time.Second
}

func millis(v int32) time.Duration {
	return time.Duration(v) * time.Millisecond
}
