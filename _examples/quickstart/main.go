// quickstart 演示最简单的用法：一行代码启动应用。
//
// 使用 [bootstrap.RunApp] 密封流程，自动处理：
//   - 信号监听（SIGINT/SIGTERM）
//   - 从 YAML 文件加载配置
//   - 引导创建应用
//   - 运行并优雅退出
//
// 启动：
//
//	go run ./_examples/quickstart/
//	go run ./_examples/quickstart/ -c /path/to/config.yaml
//
// 测试：
//
//	curl http://localhost:8080
//	curl http://localhost:8080/health
package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	bootstrap "github.com/tx7do/go-wind-bootstrap"

	// 适配层子模块，通过空导入自注册 builder。
	_ "github.com/tx7do/go-wind-bootstrap/log/zap"
	httpAdapter "github.com/tx7do/go-wind-bootstrap/transport/http"
)

func main() {
	// --- 配置路由（必须在 RunApp 之前）---
	httpAdapter.RegisterServerSetup(func(srv *httpAdapter.Server) {
		srv.GET("/", func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprintf(w, "Hello from go-wind-bootstrap quickstart!\n")
		})
		srv.GET("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	// --- 一行启动 ---
	if err := bootstrap.RunAppWithFlags(bootstrap.NewCommandFlags()); err != nil {
		slog.Error("application exited with error", "error", err)
		os.Exit(1)
	}
}
