// Package gcpubsub provides a bootstrap broker builder for Google Cloud Pub/Sub.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/broker/gcpubsub"
package gcpubsub

import (
	"context"
	"fmt"

	gcpubsubPlugin "github.com/tx7do/go-wind-plugins/transport/gcpubsub"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterBrokerBuilder(bootstrap.BrokerTypeGCPubSub, newBuilder)
}

func newBuilder(ctx context.Context, cfg *v1.Broker) (any, func(), error) {
	c := cfg.GetGcpubsub()
	if c == nil {
		return nil, nil, fmt.Errorf("gcpubsub: config is nil")
	}

	var opts []gcpubsubPlugin.ServerOption

	if projectID := c.GetProjectId(); projectID != "" {
		opts = append(opts, gcpubsubPlugin.WithProjectID(projectID))
	}
	if credFile := c.GetCredentialsFile(); credFile != "" {
		opts = append(opts, gcpubsubPlugin.WithCredentialsFile(credFile))
	}
	if endpoint := c.GetEndpoint(); endpoint != "" {
		opts = append(opts, gcpubsubPlugin.WithEndpoint(endpoint))
	}

	srv := gcpubsubPlugin.NewServer(opts...)
	return srv, func() {}, nil
}
