// Package asynq provides a bootstrap server builder for Asynq (Redis-based task queue).
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/transport/asynq"
//
// The builder maps the declarative [v1.Server_Asynq] config onto the plugin's
// option surface: redis_uri wins over the discrete redis_* fields; durations
// are declared in seconds; every field is optional and only applied when set.
package asynq

import (
	"fmt"
	"time"

	asynqLib "github.com/hibiken/asynq"

	asynqPlugin "github.com/tx7do/go-wind-plugins/transport/asynq"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	"github.com/tx7do/go-wind/transport"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeAsynq, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetAsynq()
	if c == nil {
		return nil, fmt.Errorf("asynq: config is nil")
	}

	var opts []asynqPlugin.Option

	// Redis 连接：URI 优先，其次离散字段。
	if uri := c.GetRedisUri(); uri != "" {
		opts = append(opts, asynqPlugin.WithRedisURI(uri))
	} else if addr := c.GetRedisAddress(); addr != "" {
		connOpt := asynqLib.RedisClientOpt{
			Addr:     addr,
			Username: c.GetRedisUsername(),
			Password: c.GetRedisPassword(),
			DB:       int(c.GetRedisDb()),
		}
		if poolSize := c.GetRedisPoolSize(); poolSize > 0 {
			connOpt.PoolSize = int(poolSize)
		}
		opts = append(opts, asynqPlugin.WithRedisConnOpt(connOpt))
	}

	if codec := c.GetCodec(); codec != "" {
		opts = append(opts, asynqPlugin.WithCodec(codec))
	}
	if concurrency := c.GetConcurrency(); concurrency > 0 {
		opts = append(opts, asynqPlugin.WithConcurrency(concurrency))
	}
	if location := c.GetLocation(); location != "" {
		opts = append(opts, asynqPlugin.WithLocation(location))
	}
	if len(c.GetQueues()) > 0 {
		opts = append(opts, asynqPlugin.WithQueues(c.GetQueues()))
	}
	opts = append(opts,
		asynqPlugin.WithGracefullyShutdown(c.GetEnableGracefullyShutdown()),
		asynqPlugin.WithStrictPriority(c.GetEnableStrictPriority()),
	)

	if v := c.GetShutdownTimeoutSeconds(); v > 0 {
		opts = append(opts, asynqPlugin.WithShutdownTimeout(time.Duration(v)*time.Second))
	}
	if v := c.GetHealthCheckIntervalSeconds(); v > 0 {
		opts = append(opts, asynqPlugin.WithHealthCheckInterval(time.Duration(v)*time.Second))
	}
	if v := c.GetTaskCheckIntervalSeconds(); v > 0 {
		opts = append(opts, asynqPlugin.WithTaskCheckInterval(time.Duration(v)*time.Second))
	}
	if v := c.GetDelayedTaskCheckIntervalSeconds(); v > 0 {
		opts = append(opts, asynqPlugin.WithDelayedTaskCheckInterval(time.Duration(v)*time.Second))
	}
	if v := c.GetGroupGracePeriodSeconds(); v > 0 {
		opts = append(opts, asynqPlugin.WithGroupGracePeriod(time.Duration(v)*time.Second))
	}
	if v := c.GetGroupMaxDelaySeconds(); v > 0 {
		opts = append(opts, asynqPlugin.WithGroupMaxDelay(time.Duration(v)*time.Second))
	}
	if v := c.GetGroupMaxSize(); v > 0 {
		opts = append(opts, asynqPlugin.WithGroupMaxSize(v))
	}
	if v := c.GetJanitorIntervalSeconds(); v > 0 {
		opts = append(opts, asynqPlugin.WithJanitorInterval(time.Duration(v)*time.Second))
	}
	if v := c.GetJanitorBatchSize(); v > 0 {
		opts = append(opts, asynqPlugin.WithJanitorBatchSize(v))
	}

	srv := asynqPlugin.NewServer(opts...)
	return srv, nil
}
