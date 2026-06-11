package bootstrap

import (
	"context"
	"fmt"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

// resolveConfig 检查 Config 配置中每个 optional 字段，
// 对已设置的配置源类型分别调用对应 action。
func resolveConfig(ctx context.Context, cfg *v1.Config) (func(), error) {
	type field struct {
		name   string
		action ConfigAction
	}

	var fields []field

	if cfg.GetFile() != nil {
		a, err := getConfigAction(ConfigTypeFile)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field{name: ConfigTypeFile, action: a})
	}
	if cfg.GetFs() != nil {
		a, err := getConfigAction(ConfigTypeFs)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field{name: ConfigTypeFs, action: a})
	}
	if cfg.GetEtcd() != nil {
		a, err := getConfigAction(ConfigTypeEtcd)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field{name: ConfigTypeEtcd, action: a})
	}
	if cfg.GetNacos() != nil {
		a, err := getConfigAction(ConfigTypeNacos)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field{name: ConfigTypeNacos, action: a})
	}
	if cfg.GetConsul() != nil {
		a, err := getConfigAction(ConfigTypeConsul)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field{name: ConfigTypeConsul, action: a})
	}
	if cfg.GetApollo() != nil {
		a, err := getConfigAction(ConfigTypeApollo)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field{name: ConfigTypeApollo, action: a})
	}
	if cfg.GetKubernetes() != nil {
		a, err := getConfigAction(ConfigTypeKubernetes)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field{name: ConfigTypeKubernetes, action: a})
	}
	if cfg.GetRedis() != nil {
		a, err := getConfigAction(ConfigTypeRedis)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field{name: ConfigTypeRedis, action: a})
	}
	if cfg.GetZookeeper() != nil {
		a, err := getConfigAction(ConfigTypeZookeeper)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field{name: ConfigTypeZookeeper, action: a})
	}
	if cfg.GetVault() != nil {
		a, err := getConfigAction(ConfigTypeVault)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field{name: ConfigTypeVault, action: a})
	}
	if cfg.GetHttp() != nil {
		a, err := getConfigAction(ConfigTypeHTTP)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field{name: ConfigTypeHTTP, action: a})
	}
	if cfg.GetEnv() != nil {
		a, err := getConfigAction(ConfigTypeEnv)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field{name: ConfigTypeEnv, action: a})
	}
	if cfg.GetOss() != nil {
		a, err := getConfigAction(ConfigTypeOSS)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field{name: ConfigTypeOSS, action: a})
	}
	if cfg.GetPolaris() != nil {
		a, err := getConfigAction(ConfigTypePolaris)
		if err != nil {
			return nil, err
		}
		fields = append(fields, field{name: ConfigTypePolaris, action: a})
	}

	if len(fields) == 0 {
		return nil, fmt.Errorf("bootstrap: no config source specified")
	}

	var cleanups []func()
	for _, f := range fields {
		cleanup, err := f.action(ctx, cfg)
		if err != nil {
			for _, c := range cleanups {
				c()
			}
			return nil, fmt.Errorf("bootstrap: build config %q: %w", f.name, err)
		}
		if cleanup != nil {
			cleanups = append(cleanups, cleanup)
		}
	}

	cleanup := func() {
		for _, c := range cleanups {
			c()
		}
	}
	return cleanup, nil
}
