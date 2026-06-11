// Package s3 provides a bootstrap storage builder for S3-compatible object storage.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/oss/s3"
package s3

import (
	"context"
	"fmt"

	s3Plugin "github.com/tx7do/go-wind-plugins/oss/s3"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterStorageBuilder(bootstrap.StorageTypeS3, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Storage) (any, func(), error) {
	c := cfg.GetS3()
	if c == nil {
		return nil, nil, fmt.Errorf("s3: config is nil")
	}

	pluginCfg := &s3Plugin.Config{
		Endpoint:       c.GetEndpoint(),
		Region:         c.GetRegion(),
		AccessKey:      c.GetAccessKey(),
		SecretKey:      c.GetSecretKey(),
		Token:          c.GetToken(),
		UseSsl:         c.GetUseSsl(),
		ForcePathStyle: c.GetForcePathStyle(),
		Bucket:         c.GetBucket(),
	}

	client := s3Plugin.NewClient(pluginCfg)
	if client == nil {
		return nil, nil, fmt.Errorf("s3: failed to create client")
	}

	return client, func() {}, nil
}
