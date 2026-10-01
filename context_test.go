package bootstrap

import (
	"testing"

	"github.com/tx7do/go-wind/log"
)

// TestContextServerAccessor 验证 Server/Servers 访问器的实例保留与 nil 安全。
func TestContextServerAccessor(t *testing.T) {
	if got := (*Context)(nil).Server("http"); got != nil {
		t.Fatalf("nil context Server: want nil, got %v", got)
	}
	if got := (*Context)(nil).Servers(); got != nil {
		t.Fatalf("nil context Servers: want nil, got %v", got)
	}

	servers := map[string]any{"http": "fake-instance"}
	c := newContext(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, servers, nil, func() {}, nil)

	if got := c.Server("http"); got != "fake-instance" {
		t.Fatalf("Server(http): want fake-instance, got %v", got)
	}
	if got := c.Server("nope"); got != nil {
		t.Fatalf("Server(nope): want nil, got %v", got)
	}
	if got := c.Servers()["http"]; got != "fake-instance" {
		t.Fatalf("Servers()[http]: want fake-instance, got %v", got)
	}
}

// TestContextLoggerAccessor 验证 Logger/NewModuleLogger 的返回与 nil 安全。
// log.GetLogger() 返回全局默认的 nop logger，零副作用。
func TestContextLoggerAccessor(t *testing.T) {
	if got := (*Context)(nil).Logger(); got != nil {
		t.Fatalf("nil context Logger: want nil, got %v", got)
	}
	if got := (*Context)(nil).NewModuleLogger("x"); got != nil {
		t.Fatalf("nil context NewModuleLogger: want nil, got %v", got)
	}

	l := log.GetLogger()
	c := newContext(nil, nil, l, nil, nil, nil, nil, nil, nil, nil, nil, nil, func() {}, nil)

	if got := c.Logger(); got != l {
		t.Fatalf("Logger: want the resolved logger, got %v", got)
	}
	if got := c.NewModuleLogger("data/test"); got == nil {
		t.Fatalf("NewModuleLogger: want non-nil child logger, got nil")
	}
}
