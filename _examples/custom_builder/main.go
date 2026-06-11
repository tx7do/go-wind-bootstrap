// custom_builder 演示如何注册自定义 Builder 并使用完整配置启动应用。
//
// 本示例混合使用适配层和自定义 Builder：
//   - HTTP Server、gRPC Server、Zap Logger → 适配层子模块（空导入自注册）
//   - Registry、Tracer、Metrics、Broker → 自定义 stub Builder（手动注册）
//
// 启动：
//
//	go run ./_examples/custom_builder/
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
	bootstrapV1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"

	// 适配层子模块，通过空导入自注册 HTTP/gRPC server 和 Zap logger builder。
	_ "github.com/tx7do/go-wind-bootstrap/log/zap"
	grpcAdapter "github.com/tx7do/go-wind-bootstrap/transport/grpc"
	httpAdapter "github.com/tx7do/go-wind-bootstrap/transport/http"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// --- 1. 注册自定义 Builder（适配层已通过 init 自动注册）---
	registerCustomBuilders()

	// --- 2. 列出已注册的 Builder（诊断用）---
	slog.Info("registered server builders", "types", bootstrap.ListServerBuilders())
	slog.Info("registered log builders", "types", bootstrap.ListLogBuilders())

	// --- 3. 配置路由（HTTP）---
	httpAdapter.RegisterServerSetup(func(srv *httpAdapter.Server) {
		srv.GET("/", func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprintf(w, "Hello from custom_builder example!\n")
		})
		srv.GET("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	// --- 4. 从 YAML 文件加载配置 ---
	cfg, err := bootstrap.LoadConfigFromFile("config.yaml")
	if err != nil {
		slog.Error("load config failed", "error", err)
		os.Exit(1)
	}
	slog.Info("config loaded",
		"id", cfg.GetApp().GetId(),
		"version", cfg.GetApp().GetVersion(),
		"env", cfg.GetApp().GetEnv(),
	)

	// --- 5. Bootstrap ---
	app, brokers, _, _, _, _, _, _, cleanup, err := bootstrap.Bootstrap(ctx, cfg)
	if err != nil {
		slog.Error("bootstrap failed", "error", err)
		os.Exit(1)
	}
	defer cleanup()

	// --- 5b. 获取 broker 实例 ---
	if b := brokers[bootstrap.BrokerTypeKafka]; b != nil {
		slog.Info("got kafka broker instance", "broker", b)
	}

	slog.Info("application started", "id", cfg.GetApp().GetId())

	// --- 6. Run ---
	if err := app.Run(ctx); err != nil {
		slog.Error("application exited with error", "error", err)
	}
}

// registerCustomBuilders 注册自定义的 stub Builder。
// 适配层（HTTP、gRPC、Zap）已通过 init() 自动注册，这里只注册没有适配层的组件。
func registerCustomBuilders() {
	bootstrap.MustRegisterRegistryAction(bootstrap.RegistryTypeConsul, newConsulRegistry)
	bootstrap.MustRegisterTracerBuilder(bootstrap.TracerTypeOTLP, newOTLPTracer)
	bootstrap.MustRegisterMetricsBuilder(bootstrap.MetricsTypePrometheus, newPrometheusMetrics)
	bootstrap.MustRegisterBrokerBuilder(bootstrap.BrokerTypeKafka, newKafkaBroker)
}

// ---------------------------------------------------------------------------
// Stub Builders（需要外部服务，仅打印配置）
// ---------------------------------------------------------------------------

func newConsulRegistry(ctx context.Context, appCfg *bootstrapV1.App, endpoints []string, cfg *bootstrapV1.Registry) (func(), error) {
	slog.Info("building Consul registry", "app_id", appCfg.GetId(), "app_name", appCfg.GetName())
	return func() { slog.Info("Consul registry cleaned up") }, nil
}

func newOTLPTracer(cfg *bootstrapV1.Tracer) (interface{}, func(), error) {
	otlpCfg := cfg.GetOtlp()
	slog.Info("building OTLP tracer", "endpoint", otlpCfg.GetEndpoint(), "insecure", otlpCfg.GetInsecure())
	return nil, func() { slog.Info("OTLP tracer cleaned up") }, nil
}

func newPrometheusMetrics(cfg *bootstrapV1.Metrics) (func(), error) {
	promCfg := cfg.GetPrometheus()
	slog.Info("building Prometheus metrics", "addr", promCfg.GetAddr(), "path", promCfg.GetPath())
	return func() { slog.Info("Prometheus metrics cleaned up") }, nil
}

func newKafkaBroker(ctx context.Context, cfg *bootstrapV1.Broker) (any, func(), error) {
	kafkaCfg := cfg.GetKafka()
	slog.Info("building Kafka broker", "brokers", kafkaCfg.GetBrokers(), "group_id", kafkaCfg.GetGroupId())
	// 这里返回一个 stub broker 实例。实际使用中应返回真实的 broker 对象。
	return struct{ Name string }{Name: "kafka-stub"}, func() { slog.Info("Kafka broker cleaned up") }, nil
}

var (
	_ wind.Option
	_ = grpcAdapter.RegisterServiceRegistrar
)
