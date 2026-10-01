package bootstrap

import (
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/tx7do/go-wind/transport"
)

// waitServerEndpoints returns each server's endpoint, waiting briefly for
// servers that bind to a random port (":0") to finish listening.
//
// 服务器在 Start 内部才绑定监听器，配置为 ":0" 的服务器在绑定完成前
// Endpoint() 会一直返回端口 0 的占位地址，因此需要轮询等待。解析不出
// host:port 形态的端点（如 cron://、mock://）原样采纳。
func waitServerEndpoints(servers []transport.Server, timeout time.Duration) []string {
	deadline := time.Now().Add(timeout)

	var endpoints []string
	remaining := make([]transport.Server, len(servers))
	copy(remaining, servers)

	for len(remaining) > 0 && time.Now().Before(deadline) {
		var pending []transport.Server
		for _, srv := range remaining {
			ep := srv.Endpoint()
			if ep == "" {
				continue
			}
			if endpointBound(ep) {
				endpoints = append(endpoints, ep)
				continue
			}
			pending = append(pending, srv)
		}
		remaining = pending
		if len(remaining) > 0 {
			time.Sleep(50 * time.Millisecond)
		}
	}

	// 超时兜底：未完成绑定的按当前 Endpoint() 值返回（即占位地址）。
	for _, srv := range remaining {
		if ep := srv.Endpoint(); ep != "" {
			endpoints = append(endpoints, ep)
		}
	}
	return endpoints
}

// endpointBound reports whether the endpoint carries a usable (non-zero)
// port. Unparseable shapes are treated as ready.
func endpointBound(ep string) bool {
	u, err := url.Parse(ep)
	if err != nil || u.Host == "" {
		return true
	}
	_, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return true
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		return true
	}
	return p != 0
}
