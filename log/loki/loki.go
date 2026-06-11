// Package loki provides a bootstrap log builder for the Grafana Loki logging service.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/log/loki"
package loki

import (
	"fmt"
	"time"

	lokiPlugin "github.com/tx7do/go-wind-plugins/log/loki"
	windLog "github.com/tx7do/go-wind/log"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterLogBuilder(bootstrap.LoggerTypeLoki, newBuilder)
}

func newBuilder(cfg *v1.Logger) (windLog.Logger, func(), error) {
	c := cfg.GetLoki()
	if c == nil {
		return nil, nil, fmt.Errorf("loki: config is nil")
	}
	if c.GetEndpoint() == "" {
		return nil, nil, fmt.Errorf("loki: endpoint is required")
	}

	var opts []lokiPlugin.Option
	opts = append(opts, lokiPlugin.WithEndpoint(c.GetEndpoint()))

	for k, v := range c.GetLabels() {
		opts = append(opts, lokiPlugin.WithLabel(k, v))
	}
	if c.GetBatchSize() > 0 {
		opts = append(opts, lokiPlugin.WithBatchSize(int(c.GetBatchSize())))
	}
	if c.GetFlushInterval() > 0 {
		opts = append(opts, lokiPlugin.WithFlushInterval(time.Duration(c.GetFlushInterval())*time.Millisecond))
	}

	logger, err := lokiPlugin.NewLogger(opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("loki: %w", err)
	}

	cleanup := func() { logger.Close() }
	return logger, cleanup, nil
}
