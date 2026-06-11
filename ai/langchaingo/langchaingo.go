// Package langchaingo provides a bootstrap AI builder for LangChainGo models.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/ai/langchaingo"
package langchaingo

import (
	"context"
	"fmt"

	langchaingoPlugin "github.com/tx7do/go-wind-plugins/ai/langchaingo"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterAiBuilder(bootstrap.AiTypeLangChainGo, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Ai) (any, func(), error) {
	c := cfg.GetLangchaingo()
	if c == nil {
		return nil, nil, fmt.Errorf("langchaingo: config is nil")
	}

	pluginCfg := &langchaingoPlugin.Config{
		Type:      langchaingoPlugin.ModelType(c.GetModelType()),
		ModelName: c.GetModelName(),
	}

	if ts := c.GetTimeoutSeconds(); ts > 0 {
		pluginCfg.TimeoutSeconds = ts
	}

	if cloud := c.GetCloud(); cloud != nil {
		pluginCfg.Cloud = &langchaingoPlugin.CloudConfig{
			ApiKey:  cloud.GetApiKey(),
			BaseUrl: cloud.GetBaseUrl(),
		}
	}

	if local := c.GetLocal(); local != nil {
		pluginCfg.Local = &langchaingoPlugin.LocalConfig{
			Host: local.GetHost(),
			Port: local.GetPort(),
		}
	}

	model, err := langchaingoPlugin.NewModel(pluginCfg)
	if err != nil {
		return nil, nil, fmt.Errorf("langchaingo: create model: %w", err)
	}

	return model, func() {}, nil
}
