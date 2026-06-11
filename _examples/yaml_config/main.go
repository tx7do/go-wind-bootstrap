// yaml_config 演示从 YAML 配置文件启动应用。
//
// 核心流程：
//  1. 配置路由注册回调
//  2. 从 YAML 文件加载配置
//  3. 调用 Bootstrap 创建应用并运行
//
// 本示例使用适配层子模块提供的 Zap logger 和 HTTP server builder，
// 无需手写任何 builder 代码。
//
// 启动：
//
//	go run ./_examples/yaml_config/
//
// 测试：
//
//	curl http://localhost:8080
//	curl http://localhost:8080/health
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/tx7do/go-wind"

	bootstrap "github.com/tx7do/go-wind-bootstrap"

	// 适配层子模块，通过空导入自注册 builder。
	_ "github.com/tx7do/go-wind-bootstrap/log/zap"
	httpAdapter "github.com/tx7do/go-wind-bootstrap/transport/http"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// --- 1. 配置路由 ---
	httpAdapter.RegisterServerSetup(func(srv *httpAdapter.Server) {
		srv.GET("/", func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprintf(w, "Hello from yaml_config example!\n")
		})
		srv.GET("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	// --- 2. 从 YAML 文件加载配置 ---
	cfg, err := bootstrap.LoadConfigFromFile("config.yaml")
	if err != nil {
		slog.Error("load config failed", "error", err)
		os.Exit(1)
	}
	slog.Info("config loaded",
		"id", cfg.GetApp().GetId(),
		"version", cfg.GetApp().GetVersion(),
	)

	// --- 3. Bootstrap ---
	app, _, cleanup, err := bootstrap.Bootstrap(ctx, cfg)
	if err != nil {
		slog.Error("bootstrap failed", "error", err)
		os.Exit(1)
	}
	defer cleanup()

	slog.Info("application started", "id", cfg.GetApp().GetId())

	// --- 4. Run ---
	if err := app.Run(ctx); err != nil {
		slog.Error("application exited with error", "error", err)
	}
}

var _ wind.Option
