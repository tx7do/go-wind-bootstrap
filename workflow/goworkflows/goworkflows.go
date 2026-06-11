// Package goworkflows provides a bootstrap workflow builder for go-workflows.
//
// It creates the appropriate backend (sqlite by default) and then constructs
// a [goworkflowsPlugin.WorkflowClient].
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/workflow/goworkflows"
package goworkflows

import (
	"context"
	"fmt"

	"github.com/cschleiden/go-workflows/backend"
	"github.com/cschleiden/go-workflows/backend/sqlite"

	goworkflowsPlugin "github.com/tx7do/go-wind-plugins/workflow/goworkflows"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterWorkflowBuilder(bootstrap.WorkflowTypeGoWorkflows, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Workflow) (any, func(), error) {
	c := cfg.GetGoworkflows()
	if c == nil {
		return nil, nil, fmt.Errorf("goworkflows: config is nil")
	}

	backendType := c.GetBackend()
	if backendType == "" {
		backendType = "sqlite"
	}

	dsn := c.GetDsn()

	var b backend.Backend
	var err error

	switch backendType {
	case "sqlite":
		if dsn == "" {
			dsn = ":memory:"
		}
		b = sqlite.NewSqliteBackend(dsn)
	default:
		return nil, nil, fmt.Errorf("goworkflows: unsupported backend %q (supported: sqlite)", backendType)
	}

	client, err := goworkflowsPlugin.NewClient(b)
	if err != nil {
		return nil, nil, fmt.Errorf("goworkflows: create client: %w", err)
	}

	cleanup := func() {
		_ = client.Close()
	}

	return client, cleanup, nil
}
