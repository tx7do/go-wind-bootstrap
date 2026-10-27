// Package aliyun provides a bootstrap log builder for the Alibaba Cloud SLS logging service.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/log/aliyun"
package aliyun

import (
	"fmt"

	aliyunPlugin "github.com/tx7do/go-wind-plugins/log/aliyun"
	windLog "github.com/tx7do/go-wind/log"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterLogBuilder(bootstrap.LoggerTypeAliyun, newBuilder)
}

func newBuilder(cfg *v1.Logger) (windLog.Logger, func(), error) {
	c := cfg.GetAliyun()
	if c == nil {
		return nil, nil, fmt.Errorf("aliyun: config is nil")
	}

	// 常规配置字段 → 选项的映射由 protoc-gen-fieldmap 生成（mapping.spec.yaml）。
	var opts []aliyunPlugin.Option
	opts = MapAliyunOptions(c, opts)

	logger, err := aliyunPlugin.NewAliyunLogger(opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("aliyun: %w", err)
	}

	cleanup := func() { logger.Close() }
	return logger, cleanup, nil
}
