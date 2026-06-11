package bootstrap

import (
	"context"
	"fmt"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

// resolveStorage 检查 Storage 配置中每个 optional 字段，
// 对已设置的对象存储类型分别调用对应 builder。
// 返回按类型名索引的 storage 实例映射和统一的 cleanup 函数。
func resolveStorage(ctx context.Context, cfg *v1.Storage) (map[string]any, func(), error) {
	type field struct {
		name    string
		builder StorageBuilder
	}

	var fields []field

	if cfg.GetMinio() != nil {
		b, err := getStorageBuilder(StorageTypeMinio)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: StorageTypeMinio, builder: b})
	}
	if cfg.GetS3() != nil {
		b, err := getStorageBuilder(StorageTypeS3)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: StorageTypeS3, builder: b})
	}

	if len(fields) == 0 {
		return nil, nil, fmt.Errorf("bootstrap: no storage specified")
	}

	instances := make(map[string]any)
	var cleanups []func()
	for _, f := range fields {
		inst, cleanup, err := f.builder(ctx, cfg)
		if err != nil {
			for _, c := range cleanups {
				c()
			}
			return nil, nil, fmt.Errorf("bootstrap: build storage %q: %w", f.name, err)
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
