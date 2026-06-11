// Package zerolog provides a bootstrap log builder for the Zerolog logging library.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/log/zerolog"
package zerolog

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"

	zerologPlugin "github.com/tx7do/go-wind-plugins/log/zerolog"
	windLog "github.com/tx7do/go-wind/log"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterLogBuilder(bootstrap.LoggerTypeZerolog, newBuilder)
}

func newBuilder(cfg *v1.Logger) (windLog.Logger, func(), error) {
	c := cfg.GetZerolog()
	if c == nil {
		return nil, nil, fmt.Errorf("zerolog: config is nil")
	}

	// Level.
	zl := zerolog.InfoLevel
	switch c.GetLevel() {
	case "debug":
		zl = zerolog.DebugLevel
	case "warn":
		zl = zerolog.WarnLevel
	case "error":
		zl = zerolog.ErrorLevel
	}

	// Writer.
	var writer interface {
		Write(p []byte) (n int, err error)
	}
	switch c.GetOutputPath() {
	case "", "stdout":
		writer = zerolog.SyncWriter(os.Stdout)
	case "stderr":
		writer = zerolog.SyncWriter(os.Stderr)
	default:
		f, err := os.OpenFile(c.GetOutputPath(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, nil, fmt.Errorf("zerolog: open %s: %w", c.GetOutputPath(), err)
		}
		writer = f
	}

	// Format.
	switch c.GetFormat() {
	case "console":
		writer = zerolog.ConsoleWriter{Out: writer, TimeFormat: time.RFC3339}
	}

	zlog := zerolog.New(writer).Level(zl).With().Timestamp().Logger()
	logger := zerologPlugin.NewZerologLogger(&zlog)

	cleanup := func() {}
	if closer, ok := writer.(interface{ Close() error }); ok {
		cleanup = func() { _ = closer.Close() }
	}

	return logger, cleanup, nil
}
