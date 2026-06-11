package bootstrap

import (
	"fmt"

	"github.com/tx7do/go-wind/transport"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

// resolveServer 检查 Server 配置中每个 optional 字段，
// 对已设置的传输层类型分别调用对应 builder 创建实例。
func resolveServer(cfg *v1.Server) ([]transport.Server, func(), error) {
	var servers []transport.Server

	// 按字段依次检查，调用对应 builder。
	type fieldBuilder struct {
		name string
		typ  string
		set  bool
	}
	fields := []fieldBuilder{
		{"http", serverTypeHTTP, cfg.GetHttp() != nil},
		{"http3", serverTypeHTTP3, cfg.GetHttp3() != nil},
		{"grpc", serverTypeGRPC, cfg.GetGrpc() != nil},
		{"graphql", serverTypeGraphQL, cfg.GetGraphql() != nil},
		{"sse", serverTypeSSE, cfg.GetSse() != nil},
		{"websocket", serverTypeWebsocket, cfg.GetWebsocket() != nil},
		{"tcp", serverTypeTCP, cfg.GetTcp() != nil},
		{"udp", serverTypeUDP, cfg.GetUdp() != nil},
		{"kcp", serverTypeKCP, cfg.GetKcp() != nil},
		{"thrift", serverTypeThrift, cfg.GetThrift() != nil},
		{"trpc", serverTypeTRPC, cfg.GetTrpc() != nil},
		{"webtransport", serverTypeWebtransport, cfg.GetWebtransport() != nil},
		{"cron", serverTypeCron, cfg.GetCron() != nil},
		{"hptimer", serverTypeHPTimer, cfg.GetHptimer() != nil},
		{"mcp", serverTypeMCP, cfg.GetMcp() != nil},
		{"signalr", serverTypeSignalR, cfg.GetSignalr() != nil},
		{"socketio", serverTypeSocketIO, cfg.GetSocketio() != nil},
		{"webrtc", serverTypeWebRTC, cfg.GetWebrtc() != nil},
		{"asynq", serverTypeAsynq, cfg.GetAsynq() != nil},
		{"machinery", serverTypeMachinery, cfg.GetMachinery() != nil},
	}

	for _, f := range fields {
		if !f.set {
			continue
		}
		b, err := getServerBuilder(f.typ)
		if err != nil {
			return nil, nil, err
		}
		srv, err := b(cfg)
		if err != nil {
			return nil, nil, fmt.Errorf("bootstrap: build server %q: %w", f.name, err)
		}
		if srv != nil {
			servers = append(servers, srv)
		}
	}

	return servers, func() {}, nil
}

// 内部常量，用于 server builder 注册的 key。
const (
	serverTypeHTTP         = "http"
	serverTypeHTTP3        = "http3"
	serverTypeGRPC         = "grpc"
	serverTypeGraphQL      = "graphql"
	serverTypeSSE          = "sse"
	serverTypeWebsocket    = "websocket"
	serverTypeTCP          = "tcp"
	serverTypeUDP          = "udp"
	serverTypeKCP          = "kcp"
	serverTypeThrift       = "thrift"
	serverTypeTRPC         = "trpc"
	serverTypeWebtransport = "webtransport"
	serverTypeCron         = "cron"
	serverTypeHPTimer      = "hptimer"
	serverTypeMCP          = "mcp"
	serverTypeSignalR      = "signalr"
	serverTypeSocketIO     = "socketio"
	serverTypeWebRTC       = "webrtc"
	serverTypeAsynq        = "asynq"
	serverTypeMachinery    = "machinery"
)
