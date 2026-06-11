// Package file provides a bootstrap config action for file-based config source.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/config/file"
package file

import (
	"context"
	"fmt"

	filePlugin "github.com/tx7do/go-wind-plugins/config/file"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterConfigAction(bootstrap.ConfigTypeFile, newAction)
}

func newAction(ctx context.Context, cfg *v1.Config) (func(), error) {
	c := cfg.GetFile()
	if c == nil {
		return nil, fmt.Errorf("file: config is nil")
	}

	var opts []filePlugin.Option

	if path := c.GetPath(); path != "" {
		opts = append(opts, filePlugin.WithPath(path))
	}
	if c.GetWatch() {
		opts = append(opts, filePlugin.WithWatch(true))
	}

	src, err := filePlugin.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("file: create source: %w", err)
	}

	cleanup := func() {
		_ = src.Close()
	}
	return cleanup, nil
}
