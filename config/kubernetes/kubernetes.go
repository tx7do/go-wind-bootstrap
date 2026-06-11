// Package kubernetes provides a bootstrap config action for Kubernetes ConfigMap config source.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/config/kubernetes"
package kubernetes

import (
	"context"
	"fmt"

	kubePlugin "github.com/tx7do/go-wind-plugins/config/kubernetes"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterConfigAction(bootstrap.ConfigTypeKubernetes, newAction)
}

func newAction(ctx context.Context, cfg *v1.Config) (func(), error) {
	c := cfg.GetKubernetes()
	if c == nil {
		return nil, fmt.Errorf("kubernetes: config is nil")
	}

	var opts []kubePlugin.Option

	if ns := c.GetNamespace(); ns != "" {
		opts = append(opts, kubePlugin.WithNamespace(ns))
	}
	if label := c.GetLabelSelector(); label != "" {
		opts = append(opts, kubePlugin.WithLabelSelector(label))
	}
	if field := c.GetFieldSelector(); field != "" {
		opts = append(opts, kubePlugin.WithFieldSelector(field))
	}
	if kubeConfig := c.GetKubeConfig(); kubeConfig != "" {
		opts = append(opts, kubePlugin.WithKubeConfig(kubeConfig))
	}
	if master := c.GetMaster(); master != "" {
		opts = append(opts, kubePlugin.WithMaster(master))
	}

	kubePlugin.New(opts...)

	return func() {}, nil
}
