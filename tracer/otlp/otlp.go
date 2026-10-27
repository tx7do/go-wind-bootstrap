// Package otlp provides a bootstrap tracer builder for the OTLP exporter.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/tracer/otlp"
package otlp

import (
	"context"
	"fmt"
	"time"

	otlpPlugin "github.com/tx7do/go-wind-plugins/tracer/otlp"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterTracerBuilder(bootstrap.TracerTypeOTLP, newBuilder)
}

func newBuilder(cfg *v1.Tracer) (interface{}, func(), error) {
	c := cfg.GetOtlp()
	if c == nil {
		return nil, nil, fmt.Errorf("otlp tracer: config is nil")
	}

	// 常规配置字段 → 选项的映射由 protoc-gen-fieldmap 生成（mapping.spec.yaml）；
	// sampler 上游暂未消费，与原手写保持一致，不映射。
	var opts []otlpPlugin.Option
	opts = MapOtlpOptions(c, opts)

	tp, err := otlpPlugin.New(opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("otlp tracer: create provider: %w", err)
	}

	return tp, func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = tp.Shutdown(shutdownCtx)
	}, nil
}
