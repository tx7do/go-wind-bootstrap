// Package cron provides a bootstrap server builder for Cron scheduler.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/transport/cron"
package cron

import (
	"fmt"

	cronPlugin "github.com/tx7do/go-wind-plugins/transport/cron"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	"github.com/tx7do/go-wind/transport"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeCron, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetCron()
	if c == nil {
		return nil, fmt.Errorf("cron: config is nil")
	}

	var opts []cronPlugin.Option

	if c.GetSeconds() {
		opts = append(opts, cronPlugin.WithSeconds(true))
	}

	srv := cronPlugin.NewServer(opts...)
	return srv, nil
}
