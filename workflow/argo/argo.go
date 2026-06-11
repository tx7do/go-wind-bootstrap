// Package argo provides a bootstrap workflow builder for Argo Workflows.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/workflow/argo"
package argo

import (
	"context"
	"fmt"

	argoPlugin "github.com/tx7do/go-wind-plugins/workflow/argo"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterWorkflowBuilder(bootstrap.WorkflowTypeArgo, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Workflow) (any, func(), error) {
	c := cfg.GetArgo()
	if c == nil {
		return nil, nil, fmt.Errorf("argo: config is nil")
	}

	clientOpts := argoPlugin.ClientOptions{
		ServerURL:          c.GetServerUrl(),
		Namespace:          c.GetNamespace(),
		Token:              c.GetToken(),
		InsecureSkipVerify: c.GetInsecureSkipVerify(),
	}

	client, err := argoPlugin.NewClient(clientOpts)
	if err != nil {
		return nil, nil, fmt.Errorf("argo: create client: %w", err)
	}

	cleanup := func() {
		_ = client.Close()
	}

	return client, cleanup, nil
}
