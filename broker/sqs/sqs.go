// Package sqs provides a bootstrap broker builder for AWS SQS.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/broker/sqs"
package sqs

import (
	"context"
	"fmt"

	sqsPlugin "github.com/tx7do/go-wind-plugins/transport/sqs"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterBrokerBuilder(bootstrap.BrokerTypeSQS, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Broker) (any, func(), error) {
	c := cfg.GetSqs()
	if c == nil {
		return nil, nil, fmt.Errorf("sqs: config is nil")
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

	srv := sqsPlugin.NewServer(opts...)
	return srv, func() {}, nil
}
