// Package zap provides a bootstrap log builder for the Zap logging library.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/log/zap"
package zap

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	zapPlugin "github.com/tx7do/go-wind-plugins/log/zap"
	windLog "github.com/tx7do/go-wind/log"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterLogBuilder(bootstrap.LoggerTypeZap, newBuilder)
}

func newBuilder(cfg *v1.Logger) (windLog.Logger, func(), error) {
	c := cfg.GetZap()
	if c == nil {
		return nil, nil, fmt.Errorf("zap: config is nil")
	}

	// Encoder.
	encCfg := zap.NewProductionEncoderConfig()
	var encoder zapcore.Encoder
	switch c.GetFormat() {
	case "json":
		encoder = zapcore.NewJSONEncoder(encCfg)
	default:
		encoder = zapcore.NewConsoleEncoder(encCfg)
	}

	// Level.
	zapLevel := zap.InfoLevel
	switch c.GetLevel() {
	case "debug":
		zapLevel = zap.DebugLevel
	case "warn":
		zapLevel = zap.WarnLevel
	case "error":
		zapLevel = zap.ErrorLevel
	}

	// Writer.
	writer := zapcore.AddSync(os.Stdout)
	if path := c.GetOutputPath(); path != "" && path != "stdout" && path != "stderr" {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, nil, fmt.Errorf("zap: open %s: %w", path, err)
		}
		writer = zapcore.AddSync(f)
	}

	core := zapcore.NewCore(encoder, writer, zapLevel)
	zlog := zap.New(core)

	logger := zapPlugin.NewZapLogger(zlog)
	cleanup := func() { _ = logger.Sync() }

	return logger, cleanup, nil
}
