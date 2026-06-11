// Package influxdb provides a bootstrap database builder for InfluxDB.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/database/influxdb"
package influxdb

import (
	"context"
	"fmt"
	"time"

	influxdbCrud "github.com/tx7do/go-crud/influxdb"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterDatabaseBuilder(bootstrap.DatabaseTypeInfluxdb, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Database) (any, func(), error) {
	c := cfg.GetInfluxdb()
	if c == nil {
		return nil, nil, fmt.Errorf("influxdb: config is nil")
	}

	var options []influxdbCrud.Option

	if host := c.GetHost(); host != "" {
		options = append(options, influxdbCrud.WithHost(host))
	}
	if token := c.GetToken(); token != "" {
		options = append(options, influxdbCrud.WithToken(token))
	}
	if db := c.GetDatabase(); db != "" {
		options = append(options, influxdbCrud.WithDatabase(db))
	}
	if org := c.GetOrganization(); org != "" {
		options = append(options, influxdbCrud.WithOrganization(org))
	}
	if scheme := c.GetAuthScheme(); scheme != "" {
		options = append(options, influxdbCrud.WithAuthScheme(scheme))
	}
	if v := c.GetWriteTimeoutSeconds(); v > 0 {
		options = append(options, influxdbCrud.WithWriteTimeout(time.Duration(v)*time.Second))
	}
	if v := c.GetQueryTimeoutSeconds(); v > 0 {
		options = append(options, influxdbCrud.WithQueryTimeout(time.Duration(v)*time.Second))
	}
	if v := c.GetMaxIdleConnections(); v > 0 {
		options = append(options, influxdbCrud.WithMaxIdleConnections(int(v)))
	}
	if v := c.GetIdleConnectionTimeoutSeconds(); v > 0 {
		options = append(options, influxdbCrud.WithIdleConnectionTimeout(time.Duration(v)*time.Second))
	}

	client, err := influxdbCrud.NewClient(options...)
	if err != nil {
		return nil, nil, fmt.Errorf("influxdb: create client failed: %w", err)
	}

	return client, func() {}, nil
}
