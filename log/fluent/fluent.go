// Package fluent provides a bootstrap log builder for the Fluentd logging service.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/log/fluent"
package fluent

import (
	"fmt"
	"time"

	fluentPlugin "github.com/tx7do/go-wind-plugins/log/fluent"
	windLog "github.com/tx7do/go-wind/log"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterLogBuilder(bootstrap.LoggerTypeFluent, newBuilder)
}

func newBuilder(cfg *v1.Logger) (windLog.Logger, func(), error) {
	c := cfg.GetFluent()
	if c == nil {
		return nil, nil, fmt.Errorf("fluent: config is nil")
	}
	if c.GetAddr() == "" {
		return nil, nil, fmt.Errorf("fluent: addr is required")
	}

	var opts []fluentPlugin.Option
	if c.GetTimeout() > 0 {
		opts = append(opts, fluentPlugin.WithTimeout(time.Duration(c.GetTimeout())*time.Millisecond))
	}
	if c.GetWriteTimeout() > 0 {
		opts = append(opts, fluentPlugin.WithWriteTimeout(time.Duration(c.GetWriteTimeout())*time.Millisecond))
	}
	if c.GetBufferLimit() > 0 {
		opts = append(opts, fluentPlugin.WithBufferLimit(int(c.GetBufferLimit())))
	}
	if c.GetRetryWait() > 0 {
		opts = append(opts, fluentPlugin.WithRetryWait(int(c.GetRetryWait())))
	}
	if c.GetMaxRetry() > 0 {
		opts = append(opts, fluentPlugin.WithMaxRetry(int(c.GetMaxRetry())))
	}
	if c.GetMaxRetryWait() > 0 {
		opts = append(opts, fluentPlugin.WithMaxRetryWait(int(c.GetMaxRetryWait())))
	}
	if c.GetTagPrefix() != "" {
		opts = append(opts, fluentPlugin.WithTagPrefix(c.GetTagPrefix()))
	}
	if c.GetAsync() {
		opts = append(opts, fluentPlugin.WithAsync(true))
	}

	logger, err := fluentPlugin.NewFluentLogger(c.GetAddr(), opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("fluent: %w", err)
	}

	cleanup := func() { logger.Close() }
	return logger, cleanup, nil
}
