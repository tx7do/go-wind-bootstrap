// ddd 演示基于 DDD 分层架构的完整示例。
//
// 分层结构：
//
//	cmd/main.go           → 入口：Bootstrap + 依赖注入
//	internal/domain/      → 领域层：Entity + Repository 接口
//	internal/data/        → 数据层：GORM 实现 Repository
//	internal/service/     → 服务层：业务逻辑
//	internal/server/      → 接口层：HTTP handler
//
// 启动：
//
//	go run ./_examples/ddd/cmd/
//
// 测试：
//
//	curl -X POST http://localhost:8080/users -H 'Content-Type: application/json' -d '{"name":"Alice","email":"alice@example.com"}'
//	curl http://localhost:8080/users/1
//	curl http://localhost:8080/users
//	curl -X PUT http://localhost:8080/users/1 -H 'Content-Type: application/json' -d '{"name":"Bob"}'
//	curl -X DELETE http://localhost:8080/users/1
package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	gormCrud "github.com/tx7do/go-crud/gorm"
	bootstrap "github.com/tx7do/go-wind-bootstrap"

	// 适配层子模块，通过空导入自注册 builder。
	_ "github.com/tx7do/go-wind-bootstrap/cache/local"
	_ "github.com/tx7do/go-wind-bootstrap/database/gorm"
	_ "github.com/tx7do/go-wind-bootstrap/log/zap"
	httpAdapter "github.com/tx7do/go-wind-bootstrap/transport/http"

	// DDD 分层
	"github.com/tx7do/go-wind-bootstrap/_examples/ddd/internal/data"
	"github.com/tx7do/go-wind-bootstrap/_examples/ddd/internal/server"
	"github.com/tx7do/go-wind-bootstrap/_examples/ddd/internal/service"
)

func main() {
	// --- 1. 从 YAML 加载配置 ---
	cfg, err := bootstrap.LoadConfigFromFile("configs/config.yaml")
	if err != nil {
		slog.Error("load config failed", "error", err)
		os.Exit(1)
	}

	// --- 2. BootstrapWithContext 启动所有基础设施 ---
	ctx, err := bootstrap.BootstrapWithContext(nil, cfg)
	if err != nil {
		slog.Error("bootstrap failed", "error", err)
		os.Exit(1)
	}
	defer ctx.Cleanup()

	// --- 3. 从 Context 获取基础设施实例（依赖注入）---

	// 获取 GORM 数据库客户端
	dbClient := ctx.Database(bootstrap.DatabaseTypeGorm)
	if dbClient == nil {
		slog.Error("gorm database not configured")
		os.Exit(1)
	}
	gormDB := dbClient.(*gormCrud.Client).DB // Client 嵌入 *gorm.DB
	slog.Info("database initialized", "dialect", gormDB.Dialector.Name())

	// 获取缓存实例（本示例仅打印信息）
	if cacheInst := ctx.Cache(bootstrap.CacheTypeLocal); cacheInst != nil {
		slog.Info("cache initialized", "type", "local")
	}

	// --- 4. 组装 DDD 分层（手动依赖注入）---

	// data 层：创建 Repository（依赖 DB）
	userRepo := data.NewUserRepository(gormDB)

	// service 层：创建 Service（依赖 Repository）
	userSvc := service.NewUserService(userRepo)

	// server 层：注册 HTTP 路由（依赖 Service）
	httpAdapter.RegisterServerSetup(func(srv *httpAdapter.Server) {
		// 健康检查
		srv.GET("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintln(w, `{"status":"ok"}`)
		})
		// 用户 CRUD 路由
		server.RegisterHTTPRoutes(srv, userSvc)
	})

	// --- 5. 运行 ---
	slog.Info("application started",
		"id", cfg.GetApp().GetId(),
		"env", cfg.GetApp().GetEnv(),
	)

	if err := ctx.App().Run(nil); err != nil {
		slog.Error("application exited with error", "error", err)
	}
}
