package bootstrap

import (
	"context"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

// resolveDatabase 检查 Database 配置中每个 optional 字段，
// 对已设置的数据库类型分别调用对应 builder。
// 返回按类型名索引的数据库客户端映射和统一的 cleanup 函数。
func resolveDatabase(ctx context.Context, cfg *v1.Database) (map[string]any, func(), error) {
	type field struct {
		name    string
		builder DatabaseBuilder
	}

	var fields []field

	if cfg.GetSql() != nil {
		// SQL 通过 GORM builder 注册，key = "gorm"
		b, err := getDatabaseBuilder(DatabaseTypeGorm)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: DatabaseTypeGorm, builder: b})
	}
	if cfg.GetMongodb() != nil {
		b, err := getDatabaseBuilder(DatabaseTypeMongodb)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: DatabaseTypeMongodb, builder: b})
	}
	if cfg.GetClickhouse() != nil {
		b, err := getDatabaseBuilder(DatabaseTypeClickhouse)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: DatabaseTypeClickhouse, builder: b})
	}
	if cfg.GetDoris() != nil {
		b, err := getDatabaseBuilder(DatabaseTypeDoris)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: DatabaseTypeDoris, builder: b})
	}
	if cfg.GetElasticsearch() != nil {
		b, err := getDatabaseBuilder(DatabaseTypeElasticsearch)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: DatabaseTypeElasticsearch, builder: b})
	}
	if cfg.GetOpensearch() != nil {
		b, err := getDatabaseBuilder(DatabaseTypeOpensearch)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: DatabaseTypeOpensearch, builder: b})
	}
	if cfg.GetInfluxdb() != nil {
		b, err := getDatabaseBuilder(DatabaseTypeInfluxdb)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: DatabaseTypeInfluxdb, builder: b})
	}
	if cfg.GetCassandra() != nil {
		b, err := getDatabaseBuilder(DatabaseTypeCassandra)
		if err != nil {
			return nil, nil, err
		}
		fields = append(fields, field{name: DatabaseTypeCassandra, builder: b})
	}

	if len(fields) == 0 {
		return nil, nil, nil
	}

	instances := make(map[string]any, len(fields))
	var cleanups []func()

	for _, f := range fields {
		inst, cleanup, err := f.builder(ctx, cfg)
		if err != nil {
			// cleanup already-created instances
			for _, c := range cleanups {
				c()
			}
			return nil, nil, err
		}
		if inst != nil {
			instances[f.name] = inst
		}
		if cleanup != nil {
			cleanups = append(cleanups, cleanup)
		}
	}

	aggregateCleanup := func() {
		for i := len(cleanups) - 1; i >= 0; i-- {
			cleanups[i]()
		}
	}

	return instances, aggregateCleanup, nil
}
