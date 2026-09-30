package sse

import (
	"testing"

	ssePlugin "github.com/tx7do/go-wind-plugins/transport/sse"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

// TestNewBuilder_InvokesHooks verifies that constructor options registered via
// RegisterServerOption and post-construction callbacks registered via
// RegisterServerSetup are both applied when the server is built.
func TestNewBuilder_InvokesHooks(t *testing.T) {
	setupCalled := false
	RegisterServerOption(ssePlugin.WithStreamIdKey("sid"))

	RegisterServerSetup(func(srv *Server) {
		setupCalled = true
		if srv == nil {
			t.Error("setup received nil server")
		}
	})

	srv, err := newBuilder(&v1.Server{
		Sse: &v1.Server_Sse{
			Addr: ":18080",
			Path: "/events",
		},
	})
	if err != nil {
		t.Fatalf("newBuilder: %v", err)
	}
	if srv == nil {
		t.Fatal("server is nil")
	}
	if !setupCalled {
		t.Error("RegisterServerSetup callback was not invoked")
	}
}

// TestNewBuilder_NilConfig verifies the nil-config guard.
func TestNewBuilder_NilConfig(t *testing.T) {
	if _, err := newBuilder(&v1.Server{}); err == nil {
		t.Error("expected an error for missing sse config")
	}
}
