package bootstrap

import (
	"context"
	"fmt"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

// resolveWorkflow 检查 Workflow 配置中每个 optional 字段，
// 对已设置的工作流引擎类型分别调用对应 builder。
// 返回按类型名索引的工作流客户端实例映射和统一的 cleanup 函数。
func resolveWorkflow(ctx context.Context, cfg *v1.Workflow) (map[string]any, func(), error) {
	type field struct {
		name    string
		builder WorkflowBuilder
	}

	var fields []field

	if cfg.GetTemporal() != nil {
		b, err := getWorkflowBuilder(WorkflowTypeTemporal)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: WorkflowTypeTemporal, builder: b})
	}
	if cfg.GetArgo() != nil {
		b, err := getWorkflowBuilder(WorkflowTypeArgo)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: WorkflowTypeArgo, builder: b})
	}
	if cfg.GetConductor() != nil {
		b, err := getWorkflowBuilder(WorkflowTypeConductor)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: WorkflowTypeConductor, builder: b})
	}
	if cfg.GetGoworkflows() != nil {
		b, err := getWorkflowBuilder(WorkflowTypeGoWorkflows)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: WorkflowTypeGoWorkflows, builder: b})
	}

	if len(fields) == 0 {
		return nil, nil, fmt.Errorf("bootstrap: no workflow specified")
	}

	instances := make(map[string]any)
	var cleanups []func()
	for _, f := range fields {
		inst, cleanup, err := f.builder(ctx, cfg)
		if err != nil {
			for _, c := range cleanups {
				c()
			}
			return nil, nil, fmt.Errorf("bootstrap: build workflow %q: %w", f.name, err)
		}
		if inst != nil {
			instances[f.name] = inst
		}
		if cleanup != nil {
			cleanups = append(cleanups, cleanup)
		}
	}

	finalCleanup := func() {
		for _, c := range cleanups {
			c()
		}
	}
	return instances, finalCleanup, nil
}
