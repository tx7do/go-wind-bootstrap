// Package oss provides a bootstrap config action for S3-compatible object storage config source.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/config/oss"
package oss

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"

	ossPlugin "github.com/tx7do/go-wind-plugins/config/oss"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterConfigAction(bootstrap.ConfigTypeOSS, newAction)
}

func newAction(ctx context.Context, cfg *v1.Config) (func(), error) {
	c := cfg.GetOss()
	if c == nil {
		return nil, fmt.Errorf("oss: config is nil")
	}

	awsCfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("oss: load AWS config: %w", err)
	}

	// endpoint 在 S3 客户端回调里装配，保留手工段。
	client := awss3.NewFromConfig(awsCfg, func(o *awss3.Options) {
		if endpoint := c.GetEndpoint(); endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
	})

	// 常规配置字段 → 选项的映射由 protoc-gen-fieldmap 生成（mapping.spec.yaml）。
	var opts []ossPlugin.Option
	opts = MapOssOptions(c, opts)

	_, err = ossPlugin.New(client, opts...)
	if err != nil {
		return nil, fmt.Errorf("oss: create source: %w", err)
	}

	return func() {}, nil
}
