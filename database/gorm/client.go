// Package gorm provides a bootstrap database builder for GORM.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/database/gorm"
package gorm

import (
	"context"
	"fmt"
	"sync"
	"time"

	gormCrud "github.com/tx7do/go-crud/gorm"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

// migrateModels 本适配层维护的迁移模型注册表。
// 业务代码通过 RegisterMigrateModel 注册 PO 模型，
// 在 GORM 客户端创建时自动通过 WithGetMigrateModels 桥接到 go-crud/gorm。
var (
	migrateModelsMu sync.RWMutex
	migrateModels   []interface{}
)

// RegisterMigrateModel 注册数据库模型用于自动迁移。
// 在 import 本适配层后，通过 init() 调用此函数注册 PO 结构体。
func RegisterMigrateModel(model interface{}) {
	if model == nil {
		return
	}
	migrateModelsMu.Lock()
	defer migrateModelsMu.Unlock()
	migrateModels = append(migrateModels, model)
}

// RegisterMigrateModels 批量注册数据库模型用于自动迁移。
func RegisterMigrateModels(models ...interface{}) {
	if len(models) == 0 {
		return
	}
	migrateModelsMu.Lock()
	defer migrateModelsMu.Unlock()
	migrateModels = append(migrateModels, models...)
}

// getMigrateModelList 返回已注册的迁移模型副本。
func getMigrateModelList() []interface{} {
	migrateModelsMu.RLock()
	defer migrateModelsMu.RUnlock()
	if len(migrateModels) == 0 {
		return nil
	}
	dup := make([]interface{}, len(migrateModels))
	copy(dup, migrateModels)
	return dup
}

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

	// 桥接全局注册的迁移模型到实例函数指针。
	// go-crud/gorm 的 initGormClient.doAutoMigrate 调用 c.getMigrateModels()，
	// 而不是 resolveMigrateModels()，因此需要显式桥接。
	// 业务代码通过 RegisterMigrateModel(&MyPO{}) 全局注册模型即可。
	options = append(options, gormCrud.WithGetMigrateModels(getMigrateModelList))

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
