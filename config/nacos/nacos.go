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

	var serverConfigs []constant.ServerConfig
	for _, addr := range addrs {
		serverConfigs = append(serverConfigs, constant.ServerConfig{ContextPath: "/nacos", IpAddr: addr})
	}

	clientConfig := constant.ClientConfig{
		NamespaceId: c.GetNamespace(),
	}

	nc := nacos_client.NacosClient{}
	_ = nc.SetClientConfig(clientConfig)
	nc.SetServerConfig(serverConfigs)

	client, err := config_client.NewConfigClient(&nc)
	if err != nil {
		return nil, fmt.Errorf("nacos: create client: %w", err)
	}

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
