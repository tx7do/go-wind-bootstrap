package asynq

import (
	"testing"

	"google.golang.org/protobuf/proto"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

// TestNewBuilder_MapsUriConfig verifies that a URI-based declarative config
// (the shape the admin backend uses) builds a server without connecting —
// asynq dials Redis only on Start.
func TestNewBuilder_MapsUriConfig(t *testing.T) {
	srv, err := newBuilder(&v1.Server{
		Asynq: &v1.Server_Asynq{
			RedisUri:               proto.String("redis://:pass@localhost:6379/1"),
			Codec:                  proto.String("json"),
			Concurrency:            proto.Int32(10),
			Location:               proto.String("Asia/Shanghai"),
			Queues:                 map[string]int32{"critical": 10, "default": 5, "low": 1},
			EnableStrictPriority:   proto.Bool(true),
			ShutdownTimeoutSeconds: proto.Int32(10),
		},
	})
	if err != nil {
		t.Fatalf("newBuilder: %v", err)
	}
	if srv == nil {
		t.Fatal("server is nil")
	}
	if srv.Endpoint() == "" {
		t.Error("endpoint is empty")
	}
}

// TestNewBuilder_MapsDiscreteRedisFields verifies the discrete redis_* field
// path (used when no URI is configured).
func TestNewBuilder_MapsDiscreteRedisFields(t *testing.T) {
	srv, err := newBuilder(&v1.Server{
		Asynq: &v1.Server_Asynq{
			RedisAddress:  "localhost:6379",
			RedisPassword: "pass",
			RedisDb:       2,
			RedisPoolSize: proto.Int32(8),
			Codec:         proto.String("json"),
		},
	})
	if err != nil {
		t.Fatalf("newBuilder: %v", err)
	}
	if srv == nil {
		t.Fatal("server is nil")
	}
}

// TestNewBuilder_NilConfig verifies the nil-config guard.
func TestNewBuilder_NilConfig(t *testing.T) {
	if _, err := newBuilder(&v1.Server{}); err == nil {
		t.Error("expected an error for missing asynq config")
	}
}
