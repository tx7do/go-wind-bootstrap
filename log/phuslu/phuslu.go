// Package phuslu provides a bootstrap log builder for the Phuslu logging library.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/log/phuslu"
package phuslu

import (
	"fmt"
	"os"

	phuslog "github.com/phuslu/log"

	phusluPlugin "github.com/tx7do/go-wind-plugins/log/phuslu"
	windLog "github.com/tx7do/go-wind/log"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterLogBuilder(bootstrap.LoggerTypePhuslu, newBuilder)
}

func newBuilder(cfg *v1.Logger) (windLog.Logger, func(), error) {
	c := cfg.GetPhuslu()
	if c == nil {
		return nil, nil, fmt.Errorf("phuslu: config is nil")
	}

	l := &phuslog.Logger{}

	// Level.
	switch c.GetLevel() {
	case "debug":
		l.Level = phuslog.DebugLevel
	case "warn":
		l.Level = phuslog.WarnLevel
	case "error":
		l.Level = phuslog.ErrorLevel
	default:
		l.Level = phuslog.InfoLevel
	}

	// Writer.
	switch c.GetOutputPath() {
	case "", "stderr":
		l.Writer = &phuslog.IOWriter{Writer: os.Stderr}
	case "stdout":
		l.Writer = &phuslog.IOWriter{Writer: os.Stdout}
	default:
		f, err := os.OpenFile(c.GetOutputPath(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, nil, fmt.Errorf("phuslu: open %s: %w", c.GetOutputPath(), err)
		}
		l.Writer = &phuslog.IOWriter{Writer: f}
		return phusluPlugin.NewLoggerWith(l), func() { _ = f.Close() }, nil
	}

	// Format.
	if c.GetFormat() == "console" {
		l.Writer = &phuslog.ConsoleWriter{Writer: os.Stderr}
	}

	return phusluPlugin.NewLoggerWith(l), func() {}, nil
}
