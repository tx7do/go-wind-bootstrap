// Package hclog provides a bootstrap log builder for the HashiCorp hclog logging library.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/log/hclog"
package hclog

import (
	"fmt"
	"os"

	hclogLib "github.com/hashicorp/go-hclog"

	hclogPlugin "github.com/tx7do/go-wind-plugins/log/hclog"
	windLog "github.com/tx7do/go-wind/log"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterLogBuilder(bootstrap.LoggerTypeHclog, newBuilder)
}

func newBuilder(cfg *v1.Logger) (windLog.Logger, func(), error) {
	c := cfg.GetHclog()
	if c == nil {
		return nil, nil, fmt.Errorf("hclog: config is nil")
	}

	opts := &hclogLib.LoggerOptions{
		Name:  c.GetName(),
		Level: hclogLib.Info,
	}

	// Level.
	switch c.GetLevel() {
	case "debug":
		opts.Level = hclogLib.Debug
	case "warn":
		opts.Level = hclogLib.Warn
	case "error":
		opts.Level = hclogLib.Error
	}

	// Format.
	opts.JSONFormat = c.GetFormat() == "json"

	// Output.
	switch c.GetOutputPath() {
	case "", "stderr":
		opts.Output = os.Stderr
	case "stdout":
		opts.Output = os.Stdout
	default:
		f, err := os.OpenFile(c.GetOutputPath(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, nil, fmt.Errorf("hclog: open %s: %w", c.GetOutputPath(), err)
		}
		opts.Output = f
		l := hclogLib.New(opts)
		return hclogPlugin.NewLoggerWith(l), func() { _ = f.Close() }, nil
	}

	// IncludeLocation.
	opts.IncludeLocation = c.GetIncludeLocation()

	// TimeFormat.
	if tf := c.GetTimeFormat(); tf != "" {
		opts.TimeFormat = tf
	}

	l := hclogLib.New(opts)
	return hclogPlugin.NewLoggerWith(l), func() {}, nil
}
