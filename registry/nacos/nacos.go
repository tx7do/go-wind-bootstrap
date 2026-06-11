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

func newAction(ctx context.Context, appCfg *v1.App, endpoints []string, cfg *v1.Registry) (func(), error) {
	c := cfg.GetNacos()
	if c == nil {
		return nil, fmt.Errorf("nacos: config is nil")
	}

	addrs := c.GetServerAddrs()
	if len(addrs) == 0 {
		return nil, fmt.Errorf("nacos: no server_addrs")
	}

	var serverConfigs []constant.ServerConfig
	for _, addr := range addrs {
		serverConfigs = append(serverConfigs, constant.ServerConfig{ContextPath: "/nacos", IpAddr: addr})
	}

	clientConfig := constant.ClientConfig{
		NamespaceId: c.GetNamespace(),
	}

	nc := nacos_client.NacosClient{}
	if err := nc.SetClientConfig(clientConfig); err != nil {
		return nil, fmt.Errorf("nacos: set client config: %w", err)
	}
	if err := nc.SetServerConfig(serverConfigs); err != nil {
		return nil, fmt.Errorf("nacos: set server config: %w", err)
	}

	client, err := naming_client.NewNamingClient(&nc)
	if err != nil {
		return nil, fmt.Errorf("nacos: create client: %w", err)
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

	reg := nacosPlugin.New(client, opts...)
	_ = reg

	return func() {}, nil
}
