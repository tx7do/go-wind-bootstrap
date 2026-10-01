// Package consul provides a bootstrap registry action for Consul service registry.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/registry/consul"
package consul

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"

	consulPlugin "github.com/tx7do/go-wind-plugins/registry/consul"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterRegistryAction(bootstrap.RegistryTypeConsul, newAction)
}

func newAction(ctx context.Context, cfg *v1.Registry) (bootstrap.Registry, func(), error) {
	c := cfg.GetConsul()
	if c == nil {
		return nil, nil, fmt.Errorf("consul: config is nil")
	}

	consulCfg := api.DefaultConfig()
	if addr := c.GetAddress(); addr != "" {
		consulCfg.Address = addr
	}
	if scheme := c.GetScheme(); scheme != "" {
		consulCfg.Scheme = scheme
	}
	if token := c.GetToken(); token != "" {
		consulCfg.Token = token
	}
	if tokenFile := c.GetTokenFile(); tokenFile != "" {
		consulCfg.TokenFile = tokenFile
	}
	if namespace := c.GetNamespace(); namespace != "" {
		consulCfg.Namespace = namespace
	}
	if partition := c.GetPartition(); partition != "" {
		consulCfg.Partition = partition
	}
	if pathPrefix := c.GetPathPrefix(); pathPrefix != "" {
		consulCfg.PathPrefix = pathPrefix
	}
	if auth := c.GetBasicAuth(); auth != nil {
		consulCfg.HttpAuth = &api.HttpBasicAuth{
			Username: auth.GetUsername(),
			Password: auth.GetPassword(),
		}
	}
	if waitTime := c.GetWaitTime(); waitTime > 0 {
		consulCfg.WaitTime = time.Duration(waitTime) * time.Second
	}
	if tlsCfg := c.GetTls(); tlsCfg != nil {
		consulCfg.TLSConfig = api.TLSConfig{
			InsecureSkipVerify: tlsCfg.GetInsecureSkipVerify(),
		}
		if f := tlsCfg.GetFile(); f != nil {
			consulCfg.TLSConfig.CertFile = f.GetCertPath()
			consulCfg.TLSConfig.KeyFile = f.GetKeyPath()
			consulCfg.TLSConfig.CAFile = f.GetCaPath()
		}
		if cf := tlsCfg.GetConfig(); cf != nil {
			consulCfg.TLSConfig.CertPEM = cf.GetCertPem()
			consulCfg.TLSConfig.KeyPEM = cf.GetKeyPem()
			consulCfg.TLSConfig.CAPem = cf.GetCaPem()
		}
	}

	client, err := api.NewClient(consulCfg)
	if err != nil {
		return nil, nil, fmt.Errorf("consul: create client: %w", err)
	}

	var opts []consulPlugin.Option
	if c.GetEnableHealthCheck() {
		opts = append(opts, consulPlugin.WithHealthCheck(true))
	}
	if interval := c.GetHealthCheckInterval(); interval > 0 {
		opts = append(opts, consulPlugin.WithHealthCheckInterval(int(interval)))
	}
	// WithTimeout 是服务发现超时；discovery_timeout 优先，向后兼容旧的 health_check_timeout 映射。
	if discoveryTimeout := c.GetDiscoveryTimeout(); discoveryTimeout > 0 {
		opts = append(opts, consulPlugin.WithTimeout(time.Duration(discoveryTimeout)*time.Second))
	} else if timeout := c.GetHealthCheckTimeout(); timeout > 0 {
		opts = append(opts, consulPlugin.WithTimeout(time.Duration(timeout)*time.Second))
	}
	if c.GetHeartbeat() {
		opts = append(opts, consulPlugin.WithHeartbeat(true))
	}
	if dc := c.GetDatacenter(); dc != "" {
		opts = append(opts, consulPlugin.WithDatacenter(consulPlugin.Datacenter(dc)))
	}
	if dereg := c.GetDeregisterCriticalServiceAfter(); dereg > 0 {
		opts = append(opts, consulPlugin.WithDeregisterCriticalServiceAfter(int(dereg)))
	}

	reg := consulPlugin.New(client, opts...)

	return reg, nil, nil
}
