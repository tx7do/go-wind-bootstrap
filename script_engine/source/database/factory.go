// Package database provides SQL database as a script source.
package database

import (
	"context"
	"fmt"

	"github.com/tx7do/go-scripts/source"
	dbSource "github.com/tx7do/go-scripts/source/database"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	scriptEngine "github.com/tx7do/go-wind-bootstrap/script_engine"
)

func init() {
	scriptEngine.MustRegisterSourceFactory(v1.Script_Source_DATABASE, NewSource)
}

// NewSource 根据配置创建 Database 脚本来源。
func NewSource(cfg *v1.Script_Source) (source.Reader, error) {
	if cfg == nil {
		return nil, fmt.Errorf("database source: config is nil")
	}

	opts := cfg.GetDatabaseOptions()
	if opts == nil {
		return nil, fmt.Errorf("database source: database_options is required")
	}

	var dbOpts []dbSource.Option

	if dsn := opts.GetDsn(); dsn != "" {
		dbOpts = append(dbOpts, dbSource.WithDSN(dsn))
	}
	if table := opts.GetTable(); table != "" {
		dbOpts = append(dbOpts, dbSource.WithTable(table))
	}
	if keyCol := opts.GetKeyColumn(); keyCol != "" {
		dbOpts = append(dbOpts, dbSource.WithKeyColumn(keyCol))
	}
	if valCol := opts.GetContentColumn(); valCol != "" {
		dbOpts = append(dbOpts, dbSource.WithValueColumn(valCol))
	}

	return dbSource.New(context.Background(), dbOpts...)
}
