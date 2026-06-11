// Package tencent provides a bootstrap log builder for the Tencent Cloud CLS logging service.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/log/tencent"
package tencent

import (
	"fmt"

	tencentPlugin "github.com/tx7do/go-wind-plugins/log/tencent"
	windLog "github.com/tx7do/go-wind/log"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterLogBuilder(bootstrap.LoggerTypeTencent, newBuilder)
}

func newBuilder(cfg *v1.Logger) (windLog.Logger, func(), error) {
	c := cfg.GetTencent()
	if c == nil {
		return nil, nil, fmt.Errorf("tencent: config is nil")
	}

	var opts []tencentPlugin.Option
	if c.GetEndpoint() != "" {
		opts = append(opts, tencentPlugin.WithEndpoint(c.GetEndpoint()))
	}
	if c.GetTopicId() != "" {
		opts = append(opts, tencentPlugin.WithTopicID(c.GetTopicId()))
	}
	if c.GetAccessKey() != "" {
		opts = append(opts, tencentPlugin.WithAccessKey(c.GetAccessKey()))
	}
	if c.GetAccessSecret() != "" {
		opts = append(opts, tencentPlugin.WithAccessSecret(c.GetAccessSecret()))
	}

	logger, err := tencentPlugin.NewTencentLogger(opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("tencent: %w", err)
	}

	cleanup := func() { logger.Close() }
	return logger, cleanup, nil
}
