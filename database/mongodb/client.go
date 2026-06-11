// Package mongodb provides a bootstrap database builder for MongoDB.
package mongodb

import (
	"context"
	"fmt"
	"time"

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

	var options []mongodbCrud.Option

	if uri := c.GetUri(); uri != "" {
		options = append(options, mongodbCrud.WithURI(uri))
	}
	if db := c.GetDatabase(); db != "" {
		options = append(options, mongodbCrud.WithDatabase(db))
	}
	if u := c.GetUsername(); u != "" && c.GetPassword() != "" {
		options = append(options, mongodbCrud.WithCredentials(u, c.GetPassword()))
	}
	if v := c.GetTimeoutSeconds(); v > 0 {
		options = append(options, mongodbCrud.WithTimeout(time.Duration(v)*time.Second))
	}
	if v := c.GetConnectTimeoutSeconds(); v > 0 {
		options = append(options, mongodbCrud.WithConnectTimeout(time.Duration(v)*time.Second))
	}
	if v := c.GetServerSelectionTimeoutSeconds(); v > 0 {
		options = append(options, mongodbCrud.WithServerSelectionTimeout(time.Duration(v)*time.Second))
	}
	if v := c.GetHeartbeatIntervalSeconds(); v > 0 {
		options = append(options, mongodbCrud.WithHeartbeatInterval(time.Duration(v)*time.Second))
	}
	if v := c.GetMaxConnIdleTimeSeconds(); v > 0 {
		options = append(options, mongodbCrud.WithMaxConnIdleTime(time.Duration(v)*time.Second))
	}

	client, err := mongodbCrud.NewClient(options...)
	if err != nil {
		return nil, nil, fmt.Errorf("mongodb: create client failed: %w", err)
	}

	cleanup := func() { client.Close() }
	return client, cleanup, nil
}
