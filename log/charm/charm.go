// Package charm provides a bootstrap log builder for the Charmbracelet logging library.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/log/charm"
package charm

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"

	charmPlugin "github.com/tx7do/go-wind-plugins/log/charm"
	windLog "github.com/tx7do/go-wind/log"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterLogBuilder(bootstrap.LoggerTypeCharm, newBuilder)
}

func newBuilder(cfg *v1.Logger) (windLog.Logger, func(), error) {
	c := cfg.GetCharm()
	if c == nil {
		return nil, nil, fmt.Errorf("charm: config is nil")
	}

	l := log.New(os.Stderr)

	// Level.
	switch c.GetLevel() {
	case "debug":
		l.SetLevel(log.DebugLevel)
	case "warn":
		l.SetLevel(log.WarnLevel)
	case "error":
		l.SetLevel(log.ErrorLevel)
	default:
		l.SetLevel(log.InfoLevel)
	}

	// Format.
	switch c.GetFormat() {
	case "json":
		l.SetFormatter(log.JSONFormatter)
	default:
		l.SetFormatter(log.TextFormatter)
	}

	// Output.
	switch c.GetOutputPath() {
	case "stdout":
		l.SetOutput(os.Stdout)
	case "", "stderr":
		// default is stderr
	default:
		f, err := os.OpenFile(c.GetOutputPath(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, nil, fmt.Errorf("charm: open %s: %w", c.GetOutputPath(), err)
		}
		l.SetOutput(f)
		return charmPlugin.NewLoggerWith(l), func() { _ = f.Close() }, nil
	}

	return charmPlugin.NewLoggerWith(l), func() {}, nil
}
