// Package mcp provides a bootstrap server builder for MCP (Model Context Protocol).
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/transport/mcp"
package mcp

import (
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/server"

	mcpPlugin "github.com/tx7do/go-wind-plugins/transport/mcp"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	"github.com/tx7do/go-wind/transport"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeMCP, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetMcp()
	if c == nil {
		return nil, fmt.Errorf("mcp: config is nil")
	}

	var opts []mcpPlugin.ServerOption

	if name := c.GetName(); name != "" {
		opts = append(opts, mcpPlugin.WithServerName(name))
	}
	if ver := c.GetVersion(); ver != "" {
		opts = append(opts, mcpPlugin.WithServerVersion(ver))
	}
	if addr := c.GetAddress(); addr != "" {
		opts = append(opts, mcpPlugin.WithMCPServeAddress(addr))
	}
	if t := c.GetServerType(); t != "" {
		opts = append(opts, mcpPlugin.WithMCPServeType(convertServerType(t)))
	}

	var mcpOpts []server.ServerOption
	if c.GetRecovery() {
		mcpOpts = append(mcpOpts, server.WithRecovery())
	}
	if c.GetLogging() {
		mcpOpts = append(mcpOpts, server.WithLogging())
	}
	if c.GetToolCapabilities() {
		mcpOpts = append(mcpOpts, server.WithToolCapabilities(true))
	}
	if c.GetPromptCapabilities() {
		mcpOpts = append(mcpOpts, server.WithPromptCapabilities(true))
	}
	if c.GetElicitation() {
		mcpOpts = append(mcpOpts, server.WithElicitation())
	}
	if c.GetRoots() {
		mcpOpts = append(mcpOpts, server.WithRoots())
	}
	if c.GetResourceCapabilities() {
		mcpOpts = append(mcpOpts, server.WithResourceCapabilities(true, true))
	}
	if c.GetEnableInputSchemaValidation() {
		mcpOpts = append(mcpOpts, server.WithInputSchemaValidation())
	}
	if c.GetEnableOutputSchemaValidation() {
		mcpOpts = append(mcpOpts, server.WithOutputSchemaValidation())
	}
	if ins := c.GetInstructions(); ins != "" {
		mcpOpts = append(mcpOpts, server.WithInstructions(ins))
	}

	if len(mcpOpts) > 0 {
		opts = append(opts, mcpPlugin.WithMCPServerOptions(mcpOpts...))
	}

	srv := mcpPlugin.NewServer(opts...)
	return srv, nil
}

// convertServerType 配置字符串 → 插件服务器类型；未知值回退 HTTP。
func convertServerType(serverType string) mcpPlugin.ServerType {
	switch strings.TrimSpace(strings.ToLower(serverType)) {
	case "sse":
		return mcpPlugin.ServerTypeSSE
	case "http":
		return mcpPlugin.ServerTypeHTTP
	case "stdio":
		return mcpPlugin.ServerTypeStdio
	case "in_process":
		return mcpPlugin.ServerTypeInProcess
	default:
		return mcpPlugin.ServerTypeHTTP
	}
}
