// Package minio provides a bootstrap storage builder for MinIO object storage.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/oss/minio"
package minio

import (
	"context"
	"fmt"

	minioPlugin "github.com/tx7do/go-wind-plugins/oss/minio"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterStorageBuilder(bootstrap.StorageTypeMinio, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Storage) (any, func(), error) {
	c := cfg.GetMinio()
	if c == nil {
		return nil, nil, fmt.Errorf("minio: config is nil")
	}

	pluginCfg := &minioPlugin.Config{
		Endpoint:  c.GetEndpoint(),
		AccessKey: c.GetAccessKey(),
		SecretKey: c.GetSecretKey(),
		Token:     c.GetToken(),
		UseSsl:    c.GetUseSsl(),
	}

	storage := minioPlugin.NewStorage(pluginCfg)
	if storage == nil {
		return nil, nil, fmt.Errorf("minio: failed to create storage")
	}

	return storage, func() {}, nil
}
