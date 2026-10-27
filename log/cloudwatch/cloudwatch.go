// Package cloudwatch provides a bootstrap log builder for the AWS CloudWatch Logs service.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/log/cloudwatch"
package cloudwatch

import (
	"context"
	"fmt"

	cloudwatchPlugin "github.com/tx7do/go-wind-plugins/log/cloudwatch"
	windLog "github.com/tx7do/go-wind/log"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterLogBuilder(bootstrap.LoggerTypeCloudWatch, newBuilder)
}

func newBuilder(cfg *v1.Logger) (windLog.Logger, func(), error) {
	c := cfg.GetCloudwatch()
	if c == nil {
		return nil, nil, fmt.Errorf("cloudwatch: config is nil")
	}

	// 常规配置字段 → 选项的映射由 protoc-gen-fieldmap 生成（mapping.spec.yaml）。
	var opts []cloudwatchPlugin.Option
	opts = MapCloudwatchOptions(c, opts)

	logger, err := cloudwatchPlugin.NewCloudWatchLogger(context.Background(), opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("cloudwatch: %w", err)
	}

	cleanup := func() { logger.Close() }
	return logger, cleanup, nil
}
