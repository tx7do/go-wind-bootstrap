// Package nacos provides a bootstrap registry action for Nacos service registry.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/registry/nacos"
package nacos

import (
	"context"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/v2/clients/nacos_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"

	nacosPlugin "github.com/tx7do/go-wind-plugins/registry/nacos"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterRegistryAction(bootstrap.RegistryTypeNacos, newAction)
}

func newAction(ctx context.Context, cfg *v1.Registry) (bootstrap.Registry, func(), error) {
	c := cfg.GetNacos()
	if c == nil {
		return nil, nil, fmt.Errorf("nacos: config is nil")
	}

	addrs := c.GetServerAddrs()
	if len(addrs) == 0 {
		return nil, nil, fmt.Errorf("nacos: no server_addrs")
	}

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

	// 客户端常规字段（namespace/endpoint/…/open_kms）由 FillNacosConfig 生成映射。
	clientConfig := FillNacosConfig(c, &constant.ClientConfig{})
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
	if err := nc.SetClientConfig(*clientConfig); err != nil {
		return nil, nil, fmt.Errorf("nacos: set client config: %w", err)
	}
	if err := nc.SetServerConfig(serverConfigs); err != nil {
		return nil, nil, fmt.Errorf("nacos: set server config: %w", err)
	}

	client, err := naming_client.NewNamingClient(&nc)
	if err != nil {
		return nil, nil, fmt.Errorf("nacos: create client: %w", err)
	}

	var opts []nacosPlugin.Option
	if group := c.GetGroup(); group != "" {
		opts = append(opts, nacosPlugin.WithGroup(group))
	}
	if cluster := c.GetClusterName(); cluster != "" {
		opts = append(opts, nacosPlugin.WithCluster(cluster))
	}
	if weight := c.GetWeight(); weight > 0 {
		opts = append(opts, nacosPlugin.WithWeight(weight))
	}
	if prefix := c.GetPrefix(); prefix != "" {
		opts = append(opts, nacosPlugin.WithPrefix(prefix))
	}
	if kind := c.GetDefaultKind(); kind != "" {
		opts = append(opts, nacosPlugin.WithDefaultKind(kind))
	}

	reg := nacosPlugin.New(client, opts...)

	return reg, nil, nil
}
