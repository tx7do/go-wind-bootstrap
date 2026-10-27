// Package mongodb provides a bootstrap database builder for MongoDB.
package mongodb

import (
	"context"
	"fmt"

	mongodbCrud "github.com/tx7do/go-crud/mongodb"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterDatabaseBuilder(bootstrap.DatabaseTypeMongodb, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Database) (any, func(), error) {
	c := cfg.GetMongodb()
	if c == nil {
		return nil, nil, fmt.Errorf("mongodb: config is nil")
	}

	// 常规配置字段 → 选项的映射由 protoc-gen-fieldmap 生成（mapping.spec.yaml）。
	var options []mongodbCrud.Option
	options = MapMongoDBOptions(c, options)

	// 用户名密码合成一项凭证选项，保留手工段。
	if u := c.GetUsername(); u != "" && c.GetPassword() != "" {
		options = append(options, mongodbCrud.WithCredentials(u, c.GetPassword()))
	}

	client, err := mongodbCrud.NewClient(options...)
	if err != nil {
		return nil, nil, fmt.Errorf("mongodb: create client failed: %w", err)
	}

	cleanup := func() { client.Close() }
	return client, cleanup, nil
}
