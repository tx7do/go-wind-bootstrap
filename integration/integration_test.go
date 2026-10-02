// Package integration 固化注册中心装配链路的端到端回归测试。
//
// 测试在检测到本地 etcd（127.0.0.1:2379）时自动运行，否则跳过：
//
//	cd integration && go test -count=1 -v ./...
//
// 覆盖的契约（依赖 go-wind v0.0.3 的生命周期钩子语义）：
//   - 实例注册发生在应用启动后（AfterStart），配置 ":0" 随机端口时
//     注册的是真实绑定端点；
//   - 优雅停机时实例在服务器停止前被注销（BeforeStop）。
package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	bootstrap "github.com/tx7do/go-wind-bootstrap"

	_ "github.com/tx7do/go-wind-bootstrap/registry/etcd"
	_ "github.com/tx7do/go-wind-bootstrap/transport/http"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

const (
	etcdEndpoint      = "127.0.0.1:2379"
	etcdDialTimeout   = 2 * time.Second
	instancePrefix    = "/microservices/it-reg-lifecycle/"
	deregisterTimeout = 10 * time.Second
)

func requireEtcd(t *testing.T) *clientv3.Client {
	t.Helper()
	conn, err := net.DialTimeout("tcp", etcdEndpoint, etcdDialTimeout)
	if err != nil {
		t.Skipf("etcd not available at %s, skipping: %v", etcdEndpoint, err)
	}
	_ = conn.Close()

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdEndpoint},
		DialTimeout: etcdDialTimeout,
	})
	if err != nil {
		t.Skipf("etcd client: %v", err)
	}
	t.Cleanup(func() { _ = cli.Close() })
	return cli
}

// instances queries the etcd prefix and returns the raw instance payloads.
func instances(t *testing.T, cli *clientv3.Client) []string {
	t.Helper()
	resp, err := cli.Get(context.Background(), instancePrefix, clientv3.WithPrefix())
	if err != nil {
		t.Fatalf("etcd get: %v", err)
	}
	raw := make([]string, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		raw = append(raw, string(kv.Value))
	}
	return raw
}

func waitFor(t *testing.T, cond func() bool, what string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for: %s", what)
}

// TestRegistryLifecycle 验证注册/注销生命周期与随机端口解析。
func TestRegistryLifecycle(t *testing.T) {
	cli := requireEtcd(t)
	ctx := context.Background()

	// 清理同名前缀，保证断言环境干净。
	if _, err := cli.Delete(ctx, instancePrefix, clientv3.WithPrefix()); err != nil {
		t.Fatalf("clean prefix: %v", err)
	}

	cfg := &v1.BootstrapConfig{
		App: &v1.App{
			Name:    "it-reg-lifecycle",
			Version: "0.0.1",
		},
		Server: &v1.Server{
			// 随机端口：注册的必须是 AfterStart 后解析出的真实端点。
			Http: &v1.Server_Http{Addr: "127.0.0.1:0"},
		},
		Registry: &v1.Registry{
			Etcd: &v1.Registry_Etcd{
				Endpoints: []string{etcdEndpoint},
				Ttl:       10,
				MaxRetry:  3,
			},
		},
	}

	app, _, _, _, _, _, _, _, _, _, _, cleanup, err := bootstrap.Bootstrap(ctx, cfg)
	if err != nil {
		t.Fatalf("bootstrap: %v", err)
	}

	runErr := make(chan error, 1)
	go func() { runErr <- app.Run(ctx) }()

	// 注册应当很快发生（AfterStart 钩子 + 端口绑定等待）。
	waitFor(t, func() bool { return len(instances(t, cli)) == 1 }, "instance registered", 10*time.Second)

	type instancePayload struct {
		ID        string   `json:"id"`
		Name      string   `json:"name"`
		Version   string   `json:"version"`
		Endpoints []string `json:"endpoints"`
	}
	var inst instancePayload
	if err := json.Unmarshal([]byte(instances(t, cli)[0]), &inst); err != nil {
		t.Fatalf("unmarshal instance: %v", err)
	}
	if inst.Name != "it-reg-lifecycle" || inst.Version != "0.0.1" {
		t.Fatalf("instance metadata mismatch: %+v", inst)
	}
	if len(inst.Endpoints) != 1 {
		t.Fatalf("expected 1 endpoint, got %+v", inst.Endpoints)
	}
	if strings.Contains(inst.Endpoints[0], ":0") {
		t.Fatalf("registered placeholder :0 endpoint: %s", inst.Endpoints[0])
	}
	t.Logf("registered: id=%s endpoint=%s", inst.ID, inst.Endpoints[0])

	// 优雅停机：实例必须在服务器停止前被摘除。
	if err := app.Stop(ctx); err != nil {
		t.Fatalf("app stop: %v", err)
	}
	if err := <-runErr; err != nil {
		t.Fatalf("run error: %v", err)
	}
	cleanup()

	waitFor(t, func() bool { return len(instances(t, cli)) == 0 }, "instance deregistered", deregisterTimeout)

	// 幂等：cleanup 可重复调用。
	cleanup()
	if got := len(instances(t, cli)); got != 0 {
		t.Fatalf("unexpected instances after idempotent cleanup: %d", got)
	}
	fmt.Println("lifecycle e2e passed")
}
