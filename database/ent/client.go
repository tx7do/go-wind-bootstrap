// Package ent provides a bootstrap helper for Ent ORM.
//
// 由于 Ent 的客户端类型是用户项目特定的泛型（每个项目的 ent schema 不同），
// 本包不使用 SPI init() 自动注册，而是提供泛型 Helper 函数，
// 由用户在业务代码中显式调用。
//
// 用法示例：
//
//	import (
//	    entAdapter "github.com/tx7do/go-wind-bootstrap/database/ent"
//	    entCrud "github.com/tx7do/go-crud/entgo"
//	    myent "my-project/ent"
//	)
//
//	func main() {
//	    cfg := &v1.BootstrapConfig{...}
//	    dbCreator := func(drv *sql.Driver) *myent.Client {
//	        return myent.NewClient(myent.Driver(drv))
//	    }
//	    client, cleanup, err := entAdapter.NewClient(cfg, dbCreator)
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    defer cleanup()
//	    // use client.Client() as *myent.Client
//	}
package ent

import (
	"fmt"
	"time"

	entSql "entgo.io/ent/dialect/sql"
	entCrud "github.com/tx7do/go-crud/entgo"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

// DbCreator 定义创建 Ent ORM 数据库客户端的函数类型。
// 用户需要在自己的项目中实现此函数，将 ent 生成的 Client 与 Driver 关联。
//
// 示例：
//
//	dbCreator := func(drv *sql.Driver) *myent.Client {
//	    return myent.NewClient(myent.Driver(drv))
//	}
type DbCreator[T entCrud.EntClientInterface] func(drv *entSql.Driver) T

// NewClient 创建 Ent ORM 数据库客户端。
//
// 它从 Bootstrap 配置的 Database.SQL 字段中读取 driver、dsn、trace、metrics 等参数，
// 创建底层 SQL Driver，然后通过用户提供的 dbCreator 回调构造项目特定的 Ent Client。
//
// 返回值：
//   - *entCrud.EntClient[T]：包装后的 Ent 客户端，可通过 .Client() 获取原始 Ent Client
//   - func()：清理函数，关闭数据库连接
//   - error：创建过程中的错误
func NewClient[T entCrud.EntClientInterface](cfg *v1.BootstrapConfig, dbCreator DbCreator[T]) (*entCrud.EntClient[T], func(), error) {
	if cfg == nil {
		return nil, nil, fmt.Errorf("ent: bootstrap config is nil")
	}

	dbCfg := cfg.GetDatabase()
	if dbCfg == nil {
		return nil, nil, fmt.Errorf("ent: database config is nil")
	}

	sqlCfg := dbCfg.GetSql()
	if sqlCfg == nil {
		return nil, nil, fmt.Errorf("ent: sql config is nil")
	}

	if dbCreator == nil {
		return nil, nil, fmt.Errorf("ent: dbCreator is nil")
	}

	// 创建底层 SQL Driver
	drv, err := entCrud.CreateDriver(
		sqlCfg.GetDriver(),
		sqlCfg.GetSource(),
		sqlCfg.GetEnableTrace(),
		sqlCfg.GetEnableMetrics(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("ent: create driver failed: %w", err)
	}

	// 通过用户回调创建项目特定的 Ent Client
	db := dbCreator(drv)

	// 包装为统一的 EntClient
	wrapperClient := entCrud.NewEntClient(db, drv)
	if wrapperClient == nil {
		return nil, nil, fmt.Errorf("ent: failed creating ent client")
	}

	// 设置连接池参数
	if sqlCfg.GetMaxIdleConnections() > 0 || sqlCfg.GetMaxOpenConnections() > 0 {
		maxIdle := int(sqlCfg.GetMaxIdleConnections())
		maxOpen := int(sqlCfg.GetMaxOpenConnections())
		var maxLifetime time.Duration
		if v := sqlCfg.GetConnectionMaxLifetimeSeconds(); v > 0 {
			maxLifetime = time.Duration(v) * time.Second
		}
		wrapperClient.SetConnectionOption(maxIdle, maxOpen, maxLifetime)
	}

	cleanup := func() { _ = wrapperClient.Close() }
	return wrapperClient, cleanup, nil
}
