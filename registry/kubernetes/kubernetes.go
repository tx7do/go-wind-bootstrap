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
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	kubernetesPlugin "github.com/tx7do/go-wind-plugins/registry/kubernetes"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterRegistryAction(bootstrap.RegistryTypeKubernetes, newAction)
}

func newAction(ctx context.Context, cfg *v1.Registry) (bootstrap.Registry, func(), error) {
	c := cfg.GetKubernetes()
	if c == nil {
		return nil, nil, fmt.Errorf("kubernetes: config is nil")
	}

	var config *rest.Config
	var err error
	switch {
	case c.InCluster != nil && c.GetInCluster():
		config, err = rest.InClusterConfig()
	case c.GetKubeconfig() != "":
		config, err = clientcmd.BuildConfigFromFlags("", c.GetKubeconfig())
	default:
		config, err = clientcmd.BuildConfigFromFlags("", "")
	}
	if err != nil {
		return nil, nil, fmt.Errorf("kubernetes: build config: %w", err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("kubernetes: create client: %w", err)
	}

	namespace := c.GetNamespace()
	if namespace == "" {
		namespace = "default"
	}

	reg := kubernetesPlugin.New(clientSet, namespace)

	return reg, nil, nil
}
