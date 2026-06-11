// Package http provides HTTP remote pull as a script source.
package http

import (
	"context"
	"fmt"

	"github.com/tx7do/go-scripts/source"
	httpSource "github.com/tx7do/go-scripts/source/http"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	scriptEngine "github.com/tx7do/go-wind-bootstrap/script_engine"
)

func init() {
	scriptEngine.MustRegisterSourceFactory(v1.Script_Source_HTTP, NewSource)
}

// NewSource 根据配置创建 HTTP 脚本来源。
func NewSource(cfg *v1.Script_Source) (source.Reader, error) {
	if cfg == nil {
		return nil, fmt.Errorf("http source: config is nil")
	}

	opts := cfg.GetHttpOptions()
	if opts == nil {
		return nil, fmt.Errorf("http source: http_options is required")
	}

	baseURL := opts.GetBaseUrl()
	if baseURL == "" {
		return nil, fmt.Errorf("http source: base_url is required")
	}

	return httpSource.New(context.Background(), httpSource.WithBaseURL(baseURL))
}
