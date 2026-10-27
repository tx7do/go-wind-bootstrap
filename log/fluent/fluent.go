// Package fluent provides a bootstrap log builder for the Fluentd logging service.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/log/fluent"
package fluent

import (
	"fmt"

	fluentPlugin "github.com/tx7do/go-wind-plugins/log/fluent"
	windLog "github.com/tx7do/go-wind/log"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterLogBuilder(bootstrap.LoggerTypeFluent, newBuilder)
}

func newBuilder(cfg *v1.Logger) (windLog.Logger, func(), error) {
	c := cfg.GetFluent()
	if c == nil {
		return nil, nil, fmt.Errorf("fluent: config is nil")
	}
	if c.GetAddr() == "" {
		return nil, nil, fmt.Errorf("fluent: addr is required")
	}

	// 常规配置字段 → 选项的映射由 protoc-gen-fieldmap 生成（mapping.spec.yaml）；
	// addr 为构造函数必填入参，保留在手工段。
	opts := MapFluentOptions(c, nil)

	logger, err := fluentPlugin.NewFluentLogger(c.GetAddr(), opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("fluent: %w", err)
	}

	cleanup := func() { logger.Close() }
	return logger, cleanup, nil
}
