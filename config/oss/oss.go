// Package oss provides a bootstrap config action for S3-compatible object storage config source.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/config/oss"
package oss

import (
	"context"
	"fmt"
	"time"

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

	client := awss3.NewFromConfig(awsCfg, func(o *awss3.Options) {
		if endpoint := c.GetEndpoint(); endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
	})

	var opts []ossPlugin.Option
	if bucket := c.GetBucket(); bucket != "" {
		opts = append(opts, ossPlugin.WithBucket(bucket))
	}
	if key := c.GetKey(); key != "" {
		opts = append(opts, ossPlugin.WithKey(key))
	}
	if pollMs := c.GetPollInterval(); pollMs > 0 {
		opts = append(opts, ossPlugin.WithPollInterval(time.Duration(pollMs)*time.Millisecond))
	}

	_, err = ossPlugin.New(client, opts...)
	if err != nil {
		return nil, fmt.Errorf("oss: create source: %w", err)
	}

	return func() {}, nil
}
