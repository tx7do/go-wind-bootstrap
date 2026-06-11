// Package cloudwatch provides a bootstrap log builder for the AWS CloudWatch Logs service.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/log/cloudwatch"
package cloudwatch

import (
	"context"
	"fmt"
	"time"

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

	var opts []cloudwatchPlugin.Option
	if c.GetRegion() != "" {
		opts = append(opts, cloudwatchPlugin.WithRegion(c.GetRegion()))
	}
	if c.GetLogGroup() != "" {
		opts = append(opts, cloudwatchPlugin.WithLogGroup(c.GetLogGroup()))
	}
	if c.GetLogStream() != "" {
		opts = append(opts, cloudwatchPlugin.WithLogStream(c.GetLogStream()))
	}
	if c.GetBatchSize() > 0 {
		opts = append(opts, cloudwatchPlugin.WithBatchSize(int(c.GetBatchSize())))
	}
	if c.GetFlushInterval() > 0 {
		opts = append(opts, cloudwatchPlugin.WithFlushInterval(time.Duration(c.GetFlushInterval())*time.Millisecond))
	}

	logger, err := cloudwatchPlugin.NewCloudWatchLogger(context.Background(), opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("cloudwatch: %w", err)
	}

	cleanup := func() { logger.Close() }
	return logger, cleanup, nil
}
