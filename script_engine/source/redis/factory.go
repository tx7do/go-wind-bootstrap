// Package redis provides Redis as a script source.
package redis

import (
	"context"
	"fmt"

	"github.com/tx7do/go-scripts/source"
	redisSource "github.com/tx7do/go-scripts/source/redis"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	scriptEngine "github.com/tx7do/go-wind-bootstrap/script_engine"
)

func init() {
	scriptEngine.MustRegisterSourceFactory(v1.Script_Source_REDIS, NewSource)
}

// NewSource 根据配置创建 Redis 脚本来源。
func NewSource(cfg *v1.Script_Source) (source.Reader, error) {
	if cfg == nil {
		return nil, fmt.Errorf("redis source: config is nil")
	}

	opts := cfg.GetRedisOptions()
	if opts == nil {
		return nil, fmt.Errorf("redis source: redis_options is required")
	}

	var redisOpts []redisSource.Option

	if addr := opts.GetAddr(); addr != "" {
		redisOpts = append(redisOpts, redisSource.WithAddr(addr))
	}
	if db := opts.GetDb(); db > 0 {
		redisOpts = append(redisOpts, redisSource.WithDB(int(db)))
	}
	if prefix := opts.GetPrefix(); prefix != "" {
		redisOpts = append(redisOpts, redisSource.WithPrefix(prefix))
	}

	return redisSource.New(context.Background(), redisOpts...)
}
