// Package s3 provides S3/compatible object storage as a script source.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/script_engine/source/s3"
package s3

import (
	"context"
	"fmt"

	"github.com/tx7do/go-scripts/source"
	s3Source "github.com/tx7do/go-scripts/source/s3"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	scriptEngine "github.com/tx7do/go-wind-bootstrap/script_engine"
)

func init() {
	scriptEngine.MustRegisterSourceFactory(v1.Script_Source_S3, NewSource)
}

// NewSource 根据配置创建 S3 脚本来源。
func NewSource(cfg *v1.Script_Source) (source.Reader, error) {
	if cfg == nil {
		return nil, fmt.Errorf("s3 source: config is nil")
	}

	opts := cfg.GetS3Options()
	if opts == nil {
		return nil, fmt.Errorf("s3 source: s3_options is required")
	}

	bucket := opts.GetBucket()
	if bucket == "" {
		return nil, fmt.Errorf("s3 source: bucket is required")
	}

	var s3Opts []s3Source.Option

	if region := opts.GetRegion(); region != "" {
		s3Opts = append(s3Opts, s3Source.WithRegion(region))
	}
	if prefix := opts.GetPrefix(); prefix != "" {
		s3Opts = append(s3Opts, s3Source.WithPrefix(prefix))
	}

	return s3Source.New(context.Background(), bucket, s3Opts...)
}
