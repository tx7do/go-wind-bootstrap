package bootstrap

import (
	"context"
	"fmt"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

// resolveAi 检查 Ai 配置中每个 optional 字段，
// 对已设置的 AI 模型类型分别调用对应 builder。
// 返回按类型名索引的 AI 客户端实例映射和统一的 cleanup 函数。
func resolveAi(ctx context.Context, cfg *v1.Ai) (map[string]any, func(), error) {
	type field struct {
		name    string
		builder AiBuilder
	}

	var fields []field

	if cfg.GetOpenai() != nil {
		b, err := getAiBuilder(AiTypeOpenAI)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: AiTypeOpenAI, builder: b})
	}
	if cfg.GetLangchaingo() != nil {
		b, err := getAiBuilder(AiTypeLangChainGo)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: AiTypeLangChainGo, builder: b})
	}
	if cfg.GetEino() != nil {
		b, err := getAiBuilder(AiTypeEino)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: AiTypeEino, builder: b})
	}

	if len(fields) == 0 {
		return nil, nil, fmt.Errorf("bootstrap: no ai specified")
	}

	instances := make(map[string]any)
	var cleanups []func()
	for _, f := range fields {
		inst, cleanup, err := f.builder(ctx, cfg)
		if err != nil {
			for _, c := range cleanups {
				c()
			}
			return nil, nil, fmt.Errorf("bootstrap: build ai %q: %w", f.name, err)
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
