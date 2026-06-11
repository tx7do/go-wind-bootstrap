// Package aliyun provides a bootstrap log builder for the Alibaba Cloud SLS logging service.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/log/aliyun"
package aliyun

import (
	"fmt"

	aliyunPlugin "github.com/tx7do/go-wind-plugins/log/aliyun"
	windLog "github.com/tx7do/go-wind/log"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterLogBuilder(bootstrap.LoggerTypeAliyun, newBuilder)
}

func newBuilder(cfg *v1.Logger) (windLog.Logger, func(), error) {
	c := cfg.GetAliyun()
	if c == nil {
		return nil, nil, fmt.Errorf("aliyun: config is nil")
	}

	var opts []aliyunPlugin.Option
	if c.GetEndpoint() != "" {
		opts = append(opts, aliyunPlugin.WithEndpoint(c.GetEndpoint()))
	}
	if c.GetProject() != "" {
		opts = append(opts, aliyunPlugin.WithProject(c.GetProject()))
	}
	if c.GetLogstore() != "" {
		opts = append(opts, aliyunPlugin.WithLogstore(c.GetLogstore()))
	}
	if c.GetAccessKey() != "" {
		opts = append(opts, aliyunPlugin.WithAccessKey(c.GetAccessKey()))
	}
	if c.GetAccessSecret() != "" {
		opts = append(opts, aliyunPlugin.WithAccessSecret(c.GetAccessSecret()))
	}
	if c.GetSecurityToken() != "" {
		opts = append(opts, aliyunPlugin.WithSecurityToken(c.GetSecurityToken()))
	}

	logger, err := aliyunPlugin.NewAliyunLogger(opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("aliyun: %w", err)
	}

	cleanup := func() { logger.Close() }
	return logger, cleanup, nil
}
