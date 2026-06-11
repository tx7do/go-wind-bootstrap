// Package sqs provides a bootstrap server builder for SQS transport.
package sqs

import (
	"fmt"

	sqsPlugin "github.com/tx7do/go-wind-plugins/transport/sqs"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	"github.com/tx7do/go-wind/transport"
)

func init() {
	bootstrap.MustRegisterServerBuilder(bootstrap.ServerTypeSQS, newBuilder)
}

func newBuilder(cfg *v1.Server) (transport.Server, error) {
	c := cfg.GetSqs()
	if c == nil {
		return nil, fmt.Errorf("sqs: config is nil")
	}

	var opts []sqsPlugin.ServerOption
	if region := c.GetRegion(); region != "" {
		opts = append(opts, sqsPlugin.WithRegion(region))
	}
	if endpoint := c.GetEndpoint(); endpoint != "" {
		opts = append(opts, sqsPlugin.WithEndpoint(endpoint))
	}
	if queueURL := c.GetQueueUrl(); queueURL != "" {
		opts = append(opts, sqsPlugin.WithQueueUrl(queueURL))
	}
	if codec := c.GetCodec(); codec != "" {
		opts = append(opts, sqsPlugin.WithCodec(codec))
	}

	return sqsPlugin.NewServer(opts...), nil
}
