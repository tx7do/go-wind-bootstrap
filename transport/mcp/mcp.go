// Package mcp provides a bootstrap server builder for MCP (Model Context Protocol).
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/transport/mcp"
package mcp

import (
	"fmt"

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

	srv := mcpPlugin.NewServer(opts...)
	return srv, nil
}
