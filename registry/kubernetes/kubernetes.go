// Package kubernetes provides a bootstrap registry action for Kubernetes service registry.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/registry/kubernetes"
package kubernetes

import (
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	kubernetesPlugin "github.com/tx7do/go-wind-plugins/registry/kubernetes"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterRegistryAction(bootstrap.RegistryTypeKubernetes, newAction)
}

func newAction(ctx context.Context, appCfg *v1.App, endpoints []string, cfg *v1.Registry) (func(), error) {
	c := cfg.GetKubernetes()
	if c == nil {
		return nil, fmt.Errorf("kubernetes: config is nil")
	}

	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		return nil, fmt.Errorf("kubernetes: build config: %w", err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("kubernetes: create client: %w", err)
	}

	namespace := c.GetNamespace()
	if namespace == "" {
		namespace = "default"
	}

	reg := kubernetesPlugin.New(clientSet, namespace)
	_ = reg

	return func() {}, nil
}
