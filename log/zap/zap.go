// Package zap provides a bootstrap log builder for the Zap logging library.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/log/zap"
package zap

import (
	"fmt"
	"os"

	"github.com/natefinch/lumberjack"
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
	switch c.GetWriter() {
	case "stdout", "":
		if c.GetWriter() == "" {
			if path := c.GetOutputPath(); path != "" && path != "stdout" && path != "stderr" {
				w, err := newFileWriter(c)
				if err != nil {
					return nil, nil, err
				}
				writer = w
			}
		}
	case "stderr":
		writer = zapcore.AddSync(os.Stderr)
	case "file":
		if path := c.GetOutputPath(); path != "" {
			w, err := newFileWriter(c)
			if err != nil {
				return nil, nil, err
			}
			writer = w
		}
	default:
		return nil, nil, fmt.Errorf("zap: unknown writer %q (want stdout/stderr/file)", c.GetWriter())
	}

	core := zapcore.NewCore(encoder, writer, zapLevel)
	zlog := zap.New(core)

	logger := zapPlugin.NewZapLogger(zlog)
	cleanup := func() { _ = logger.Sync() }

	return logger, cleanup, nil
}

// newFileWriter builds the file writer: lumberjack rolling when rolling
// params are configured, a plain append-only file otherwise.
func newFileWriter(c *v1.Logger_Zap) (zapcore.WriteSyncer, error) {
	path := c.GetOutputPath()
	if path == "" {
		return nil, fmt.Errorf("zap: writer=file requires output_path")
	}
	if c.GetMaxSizeMb() > 0 || c.GetMaxAgeDays() > 0 || c.GetMaxBackups() > 0 {
		return zapcore.AddSync(&lumberjack.Logger{
			Filename:   path,
			MaxSize:    int(c.GetMaxSizeMb()),
			MaxAge:     int(c.GetMaxAgeDays()),
			MaxBackups: int(c.GetMaxBackups()),
			Compress:   c.GetCompress(),
		}), nil
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("zap: open %s: %w", path, err)
	}
	return zapcore.AddSync(f), nil
}
