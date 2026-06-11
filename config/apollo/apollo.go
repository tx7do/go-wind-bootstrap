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

	var opts []apolloPlugin.Option

	if appID := c.GetAppId(); appID != "" {
		opts = append(opts, apolloPlugin.WithAppID(appID))
	}
	if cluster := c.GetCluster(); cluster != "" {
		opts = append(opts, apolloPlugin.WithCluster(cluster))
	}
	if endpoint := c.GetEndpoint(); endpoint != "" {
		opts = append(opts, apolloPlugin.WithEndpoint(endpoint))
	}
	if ns := c.GetNamespace(); ns != "" {
		opts = append(opts, apolloPlugin.WithNamespace(ns))
	}
	if secret := c.GetSecret(); secret != "" {
		opts = append(opts, apolloPlugin.WithSecret(secret))
	}
	if c.GetIsBackupConfig() {
		opts = append(opts, apolloPlugin.WithEnableBackup())
	} else {
		opts = append(opts, apolloPlugin.WithDisableBackup())
	}
	if backupPath := c.GetBackupPath(); backupPath != "" {
		opts = append(opts, apolloPlugin.WithBackupPath(backupPath))
	}

	apolloPlugin.NewSource(opts...)

	return func() {}, nil
}
