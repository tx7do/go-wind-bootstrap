// Package gorm provides a bootstrap database builder for GORM.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/database/gorm"
package gorm

import (
	"context"
	"fmt"
	"time"

	gormCrud "github.com/tx7do/go-crud/gorm"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterDatabaseBuilder(bootstrap.DatabaseTypeGorm, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Database) (any, func(), error) {
	c := cfg.GetSql()
	if c == nil {
		return nil, nil, fmt.Errorf("gorm: sql config is nil")
	}

	var options []gormCrud.Option

	if driver := c.GetDriver(); driver != "" {
		options = append(options, gormCrud.WithDriverName(driver))
	}
	if source := c.GetSource(); source != "" {
		options = append(options, gormCrud.WithDSN(source))
	}

	options = append(options, gormCrud.WithEnableMigrate(c.GetMigrate()))
	options = append(options, gormCrud.WithEnableTrace(c.GetEnableTrace()))
	options = append(options, gormCrud.WithEnableMetrics(c.GetEnableMetrics()))

	if v := c.GetMaxIdleConnections(); v > 0 {
		options = append(options, gormCrud.WithMaxIdleConns(int(v)))
	}
	if v := c.GetMaxOpenConnections(); v > 0 {
		options = append(options, gormCrud.WithMaxOpenConns(int(v)))
	}
	if v := c.GetConnectionMaxLifetimeSeconds(); v > 0 {
		options = append(options, gormCrud.WithConnMaxLifetime(time.Duration(v)*time.Second))
	}

	db, err := gormCrud.NewClient(options...)
	if err != nil {
		return nil, nil, fmt.Errorf("gorm: create client failed: %w", err)
	}

	cleanup := func() {
		sqlDB, _ := db.DB.DB()
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
	}
	return db, cleanup, nil
}
