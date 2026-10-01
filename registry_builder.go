package bootstrap

import (
	"context"
	"fmt"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

// resolveRegistry 检查 Registry 配置中每个 optional 字段，
// 对已设置的注册中心类型分别调用对应 action，返回按类型名键控的
// Registry 实例映射（供实例注册钩子与 Context.Registry 访问器使用）。
func resolveRegistry(ctx context.Context, cfg *v1.Registry) (map[string]Registry, func(), error) {
	type field struct {
		name   string
		action RegistryAction
	}

	var fields []field

	if cfg.GetConsul() != nil {
		a, err := getRegistryAction(RegistryTypeConsul)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: RegistryTypeConsul, action: a})
	}
	if cfg.GetEtcd() != nil {
		a, err := getRegistryAction(RegistryTypeEtcd)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: RegistryTypeEtcd, action: a})
	}
	if cfg.GetNacos() != nil {
		a, err := getRegistryAction(RegistryTypeNacos)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: RegistryTypeNacos, action: a})
	}
	if cfg.GetZookeeper() != nil {
		a, err := getRegistryAction(RegistryTypeZookeeper)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: RegistryTypeZookeeper, action: a})
	}
	if cfg.GetPolaris() != nil {
		a, err := getRegistryAction(RegistryTypePolaris)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: RegistryTypePolaris, action: a})
	}
	if cfg.GetEureka() != nil {
		a, err := getRegistryAction(RegistryTypeEureka)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: RegistryTypeEureka, action: a})
	}
	if cfg.GetKubernetes() != nil {
		a, err := getRegistryAction(RegistryTypeKubernetes)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: RegistryTypeKubernetes, action: a})
	}
	if cfg.GetServiceComb() != nil {
		a, err := getRegistryAction(RegistryTypeServiceComb)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: RegistryTypeServiceComb, action: a})
	}

	if len(fields) == 0 {
		return nil, nil, fmt.Errorf("bootstrap: no registry specified")
	}

	registries := make(map[string]Registry, len(fields))
	var cleanups []func()
	for _, f := range fields {
		reg, cleanup, err := f.action(ctx, cfg)
		if err != nil {
			for _, c := range cleanups {
				c()
			}
			return nil, nil, fmt.Errorf("bootstrap: build registry %q: %w", f.name, err)
		}
		registries[f.name] = reg
		if cleanup != nil {
			cleanups = append(cleanups, cleanup)
		}
	}

	finalCleanup := func() {
		for _, c := range cleanups {
			c()
		}
	}
	return registries, finalCleanup, nil
}
