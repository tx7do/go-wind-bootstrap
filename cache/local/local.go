// Package local provides a bootstrap cache builder for in-memory FreeCache.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/cache/local"
package local

import (
	"context"
	"fmt"
	"time"

	localPlugin "github.com/tx7do/go-wind-plugins/cache/local"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterCacheBuilder(bootstrap.CacheTypeLocal, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Cache) (any, func(), error) {
	c := cfg.GetLocal()
	if c == nil {
		return nil, nil, fmt.Errorf("local cache: config is nil")
	}

	var opts []localPlugin.Option

	if size := c.GetSize(); size > 0 {
		opts = append(opts, localPlugin.WithSize(int(size)))
	}
	if ttlSec := c.GetDefaultTtlSeconds(); ttlSec > 0 {
		opts = append(opts, localPlugin.WithDefaultTTL(time.Duration(ttlSec)*time.Second))
	}

	cache := localPlugin.New(opts...)

	cleanup := func() {
		_ = cache.Close()
	}

	return cache, cleanup, nil
}
