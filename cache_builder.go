package bootstrap

import (
	"context"
	"fmt"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

// resolveCache 检查 Cache 配置中每个 optional 字段，
// 对已设置的缓存类型分别调用对应 builder。
// 返回按类型名索引的 cache 实例映射和统一的 cleanup 函数。
func resolveCache(ctx context.Context, cfg *v1.Cache) (map[string]any, func(), error) {
	type field struct {
		name    string
		builder CacheBuilder
	}

	var fields []field

	if cfg.GetLocal() != nil {
		b, err := getCacheBuilder(CacheTypeLocal)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: CacheTypeLocal, builder: b})
	}
	if cfg.GetRedis() != nil {
		b, err := getCacheBuilder(CacheTypeRedis)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: CacheTypeRedis, builder: b})
	}

	if len(fields) == 0 {
		return nil, nil, fmt.Errorf("bootstrap: no cache specified")
	}

	instances := make(map[string]any)
	var cleanups []func()
	for _, f := range fields {
		inst, cleanup, err := f.builder(ctx, cfg)
		if err != nil {
			for _, c := range cleanups {
				c()
			}
			return nil, nil, fmt.Errorf("bootstrap: build cache %q: %w", f.name, err)
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
