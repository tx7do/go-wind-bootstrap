// Package logrus provides a bootstrap log builder for the Logrus logging library.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/log/logrus"
package logrus

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	logrusPlugin "github.com/tx7do/go-wind-plugins/log/logrus"
	windLog "github.com/tx7do/go-wind/log"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterLogBuilder(bootstrap.LoggerTypeLogrus, newBuilder)
}

func newBuilder(cfg *v1.Logger) (windLog.Logger, func(), error) {
	c := cfg.GetLogrus()
	if c == nil {
		return nil, nil, fmt.Errorf("logrus: config is nil")
	}

	l := logrus.New()

	// Level.
	level, err := logrus.ParseLevel(c.GetLevel())
	if err != nil {
		level = logrus.InfoLevel
	}
	l.SetLevel(level)

	// Format.
	switch c.GetFormat() {
	case "json":
		l.SetFormatter(&logrus.JSONFormatter{})
	default:
		l.SetFormatter(&logrus.TextFormatter{})
	}

	// Output.
	switch c.GetOutputPath() {
	case "", "stdout":
		l.SetOutput(os.Stdout)
	case "stderr":
		l.SetOutput(os.Stderr)
	default:
		f, err := os.OpenFile(c.GetOutputPath(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, nil, fmt.Errorf("logrus: open %s: %w", c.GetOutputPath(), err)
		}
		l.SetOutput(f)
		cleanup := func() { _ = f.Close() }
		return logrusPlugin.NewLogrusLogger(l), cleanup, nil
	}

	return logrusPlugin.NewLogrusLogger(l), func() {}, nil
}
