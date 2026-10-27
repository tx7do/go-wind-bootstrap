// Package nacos provides a bootstrap config action for Nacos config source.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/config/nacos"
package nacos

import (
	"context"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/nacos_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"

	nacosPlugin "github.com/tx7do/go-wind-plugins/config/nacos"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterConfigAction(bootstrap.ConfigTypeNacos, newAction)
}

func newAction(ctx context.Context, cfg *v1.Config) (func(), error) {
	c := cfg.GetNacos()
	if c == nil {
		return nil, fmt.Errorf("nacos: config is nil")
	}

	addrs := c.GetServerAddrs()
	if len(addrs) == 0 {
		return nil, fmt.Errorf("nacos: no server_addrs")
	}

	// ServerConfig 组装与多地址展开保留手工段。
	serverCfg := constant.ServerConfig{ContextPath: "/nacos"}
	if scheme := c.GetScheme(); scheme != "" {
		serverCfg.Scheme = scheme
	}
	if contextPath := c.GetContextPath(); contextPath != "" {
		serverCfg.ContextPath = contextPath
	}
	if grpcPort := c.GetGrpcPort(); grpcPort > 0 {
		serverCfg.GrpcPort = grpcPort
	}

	var serverConfigs []constant.ServerConfig
	for _, addr := range addrs {
		sc := serverCfg
		sc.IpAddr = addr
		serverConfigs = append(serverConfigs, sc)
	}

	// 常规客户端字段 → ClientConfig 的映射由 protoc-gen-fieldmap 生成（mapping.spec.yaml）。
	clientConfig := constant.ClientConfig{}
	FillNacosClientConfig(c, &clientConfig)

	// nacos SDK 的 TLS 仅支持文件路径，保留手工段。
	if tlsCfg := c.GetTls(); tlsCfg != nil {
		clientConfig.TLSCfg = constant.TLSConfig{
			Enable:   true,
			TrustAll: tlsCfg.GetInsecureSkipVerify(),
		}
		if f := tlsCfg.GetFile(); f != nil {
			clientConfig.TLSCfg.CertFile = f.GetCertPath()
			clientConfig.TLSCfg.KeyFile = f.GetKeyPath()
			clientConfig.TLSCfg.CaFile = f.GetCaPath()
		}
	}

	nc := nacos_client.NacosClient{}
	_ = nc.SetClientConfig(clientConfig)
	nc.SetServerConfig(serverConfigs)

	client, err := config_client.NewConfigClient(&nc)
	if err != nil {
		return nil, fmt.Errorf("nacos: create client: %w", err)
	}

	// group/data_id 为插件选项，保留手工段。
	var opts []nacosPlugin.Option
	if group := c.GetGroup(); group != "" {
		opts = append(opts, nacosPlugin.WithGroup(group))
	}
	if dataID := c.GetDataId(); dataID != "" {
		opts = append(opts, nacosPlugin.WithDataID(dataID))
	}

	nacosPlugin.New(client, opts...)

	return func() {}, nil
}
