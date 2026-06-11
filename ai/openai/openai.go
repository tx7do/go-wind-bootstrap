// Package openai provides a bootstrap AI builder for OpenAI-compatible clients.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/ai/openai"
package openai

import (
	"context"
	"fmt"

	openaiPlugin "github.com/tx7do/go-wind-plugins/ai/openai"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterAiBuilder(bootstrap.AiTypeOpenAI, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Ai) (any, func(), error) {
	c := cfg.GetOpenai()
	if c == nil {
		return nil, nil, fmt.Errorf("openai: config is nil")
	}

	pluginCfg := &openaiPlugin.Config{
		Type:      openaiPlugin.ModelType(c.GetModelType()),
		ModelName: c.GetModelName(),
	}

	if ts := c.GetTimeoutSeconds(); ts > 0 {
		pluginCfg.TimeoutSeconds = ts
	}

	if cloud := c.GetCloud(); cloud != nil {
		pluginCfg.Cloud = &openaiPlugin.CloudConfig{
			ApiKey:       cloud.GetApiKey(),
			BaseUrl:      cloud.GetBaseUrl(),
			Organization: cloud.GetOrganization(),
		}
	}

	if local := c.GetLocal(); local != nil {
		pluginCfg.Local = &openaiPlugin.LocalConfig{
			Host: local.GetHost(),
			Port: local.GetPort(),
		}
	}

	client, err := openaiPlugin.NewClient(pluginCfg)
	if err != nil {
		return nil, nil, fmt.Errorf("openai: create client: %w", err)
	}

	return client, func() {}, nil
}
