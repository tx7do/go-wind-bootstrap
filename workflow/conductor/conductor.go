// Package conductor provides a bootstrap workflow builder for Netflix Conductor.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/workflow/conductor"
package conductor

import (
	"context"
	"fmt"

	conductorPlugin "github.com/tx7do/go-wind-plugins/workflow/conductor"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterWorkflowBuilder(bootstrap.WorkflowTypeConductor, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Workflow) (any, func(), error) {
	c := cfg.GetConductor()
	if c == nil {
		return nil, nil, fmt.Errorf("conductor: config is nil")
	}

	clientOpts := conductorPlugin.ClientOptions{
		ServerURL:  c.GetServerUrl(),
		AuthKey:    c.GetAuthKey(),
		AuthSecret: c.GetAuthSecret(),
	}

	client, err := conductorPlugin.NewClient(clientOpts)
	if err != nil {
		return nil, nil, fmt.Errorf("conductor: create client: %w", err)
	}

	cleanup := func() {
		_ = client.Close()
	}

	return client, cleanup, nil
}
