// Package apollo provides a bootstrap config action for Apollo config source.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/config/apollo"
package apollo

import (
	"context"
	"fmt"

	apolloPlugin "github.com/tx7do/go-wind-plugins/config/apollo"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterConfigAction(bootstrap.ConfigTypeApollo, newAction)
}

func newAction(ctx context.Context, cfg *v1.Config) (func(), error) {
	c := cfg.GetApollo()
	if c == nil {
		return nil, fmt.Errorf("apollo: config is nil")
	}

	// 常规配置字段 → 选项的映射由 protoc-gen-fieldmap 生成（mapping.spec.yaml）。
	var opts []apolloPlugin.Option
	opts = MapApolloOptions(c, opts)

	// 备份开关是无参选项的布尔分支，保留手工段。
	if c.GetIsBackupConfig() {
		opts = append(opts, apolloPlugin.WithEnableBackup())
	} else {
		opts = append(opts, apolloPlugin.WithDisableBackup())
	}

	apolloPlugin.NewSource(opts...)

	return func() {}, nil
}
