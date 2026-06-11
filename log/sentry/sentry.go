// Package sentry provides a bootstrap log builder for the Sentry error tracking service.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/log/sentry"
package sentry

import (
	"fmt"

	sentryPlugin "github.com/tx7do/go-wind-plugins/log/sentry"
	windLog "github.com/tx7do/go-wind/log"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterLogBuilder(bootstrap.LoggerTypeSentry, newBuilder)
}

func newBuilder(cfg *v1.Logger) (windLog.Logger, func(), error) {
	c := cfg.GetSentry()
	if c == nil {
		return nil, nil, fmt.Errorf("sentry: config is nil")
	}
	if c.GetDsn() == "" {
		return nil, nil, fmt.Errorf("sentry: dsn is required")
	}

	var opts []sentryPlugin.Option
	opts = append(opts, sentryPlugin.WithDSN(c.GetDsn()))
	if c.GetEnvironment() != "" {
		opts = append(opts, sentryPlugin.WithEnvironment(c.GetEnvironment()))
	}
	if c.GetRelease() != "" {
		opts = append(opts, sentryPlugin.WithRelease(c.GetRelease()))
	}
	if c.GetServerName() != "" {
		opts = append(opts, sentryPlugin.WithServerName(c.GetServerName()))
	}

	logger, err := sentryPlugin.NewLogger(opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("sentry: %w", err)
	}

	cleanup := func() { logger.Close() }
	return logger, cleanup, nil
}
