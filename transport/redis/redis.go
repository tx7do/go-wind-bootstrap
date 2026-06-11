// Package redis provides a bootstrap server builder for Redis transport.
package redis

import (
	"fmt"

	redisPlugin "github.com/tx7do/go-wind-plugins/transport/redis"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	"github.com/tx7do/go-wind/transport"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeRedis, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetRedisServer()
	if c == nil {
		return nil, fmt.Errorf("redis: config is nil")
	}

	var opts []redisPlugin.ServerOption
	if addr := c.GetAddr(); addr != "" {
		opts = append(opts, redisPlugin.WithAddress(addr))
	}
	if codec := c.GetCodec(); codec != "" {
		opts = append(opts, redisPlugin.WithCodec(codec))
	}

	return redisPlugin.NewServer(opts...), nil
}
