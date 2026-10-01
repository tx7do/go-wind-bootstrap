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

	var opts []otlpPlugin.Option
	if endpoint := c.GetEndpoint(); endpoint != "" {
		opts = append(opts, otlpPlugin.WithEndpoint(endpoint))
	}
	if c.GetInsecure() {
		opts = append(opts, otlpPlugin.WithInsecure(true))
	}
	if c.GetUseHttp() {
		opts = append(opts, otlpPlugin.WithHTTP(true))
	}
	if headers := c.GetHeaders(); len(headers) > 0 {
		opts = append(opts, otlpPlugin.WithHeaders(headers))
	}
	if ratio := c.GetSampleRatio(); ratio > 0 {
		opts = append(opts, otlpPlugin.WithSampleRatio(ratio))
	}
	if ms := c.GetBatchTimeoutMs(); ms > 0 {
		opts = append(opts, otlpPlugin.WithBatchTimeout(time.Duration(ms)*time.Millisecond))
	}
	if ms := c.GetExportTimeoutMs(); ms > 0 {
		opts = append(opts, otlpPlugin.WithExportTimeout(time.Duration(ms)*time.Millisecond))
	}

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
