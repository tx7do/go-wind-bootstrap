// Package glog provides a bootstrap log builder for the Google glog logging library.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/log/glog"
package glog

import (
	"flag"
	"fmt"

	glogPlugin "github.com/tx7do/go-wind-plugins/log/glog"
	windLog "github.com/tx7do/go-wind/log"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterLogBuilder(bootstrap.LoggerTypeGlog, newBuilder)
}

func newBuilder(cfg *v1.Logger) (windLog.Logger, func(), error) {
	c := cfg.GetGlog()
	if c == nil {
		return nil, nil, fmt.Errorf("glog: config is nil")
	}

	// glog 使用命令行 flag 控制行为，设置默认值。
	if c.GetLogDir() != "" {
		_ = flag.Set("log_dir", c.GetLogDir())
	}
	if c.GetLogtostderr() {
		_ = flag.Set("logtostderr", "true")
	}
	if c.GetStderrthreshold() > 0 {
		_ = flag.Set("stderrthreshold", fmt.Sprintf("%d", c.GetStderrthreshold()))
	}
	if c.GetV() > 0 {
		_ = flag.Set("v", fmt.Sprintf("%d", c.GetV()))
	}

	// 确保在首次使用前解析 flag。
	flag.Parse()

	logger := glogPlugin.NewLogger()
	cleanup := func() { logger.Close() }

	return logger, cleanup, nil
}
