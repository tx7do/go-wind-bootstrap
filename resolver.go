package bootstrap

import (
	"time"

	wind "github.com/tx7do/go-wind"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

// resolveApp converts [App] into a slice of [wind.Option].
func resolveApp(cfg *v1.App) []wind.Option {
	var opts []wind.Option

	if cfg.GetId() != "" {
		opts = append(opts, wind.WithID(cfg.GetId()))
	}
	if cfg.GetName() != "" {
		opts = append(opts, wind.WithName(cfg.GetName()))
	}
	if cfg.GetVersion() != "" {
		opts = append(opts, wind.WithVersion(cfg.GetVersion()))
	}
	if dur := cfg.GetStopTimeout(); dur != nil {
		opts = append(opts, wind.WithStopTimeout(dur.AsDuration()))
	} else if cfg.GetEnv() == "production" {
		// 生产环境默认 30 秒优雅停机。
		opts = append(opts, wind.WithStopTimeout(30*time.Second))
	}

	return opts
}
