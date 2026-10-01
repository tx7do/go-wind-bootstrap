package bootstrap

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	wind "github.com/tx7do/go-wind"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

// InstanceRegistrar is the minimal registration surface a registry plugin
// must expose for the bootstrap layer to manage the app instance lifecycle.
// All registry plugins ([wind.Instance]-based Register/Deregister) satisfy it.
type InstanceRegistrar interface {
	Register(ctx context.Context, inst *wind.Instance) error
	Deregister(ctx context.Context, inst *wind.Instance) error
}

// RegisterInstance registers the application instance described by [appCfg]
// with the given registrar, using the (already resolved) server endpoints.
//
// Registration is skipped — returning a no-op cleanup — when the app has no
// name or there are no endpoints, i.e. when the registry acts purely as a
// client (discovery) instead of a registrar.
//
// The returned cleanup deregisters the instance; it is safe to call after the
// application context has been cancelled and uses its own short deadline.
// The caller remains responsible for closing the underlying registry client.
func RegisterInstance(ctx context.Context, reg InstanceRegistrar, appCfg *v1.App, endpoints []string) (func(), error) {
	name := appCfg.GetName()
	if name == "" || len(endpoints) == 0 {
		return func() {}, nil
	}

	inst := &wind.Instance{
		ID:        appCfg.GetId(),
		Name:      name,
		Version:   appCfg.GetVersion(),
		Endpoints: endpoints,
	}
	if inst.ID == "" {
		inst.ID = name + "-" + newInstanceID()
	}

	if err := reg.Register(ctx, inst); err != nil {
		return nil, fmt.Errorf("bootstrap: register instance %q: %w", inst.ID, err)
	}

	return func() {
		dctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = reg.Deregister(dctx, inst)
	}, nil
}

// newInstanceID returns a short random hex suffix for auto-generated
// instance IDs, avoiding a UUID dependency.
func newInstanceID() string {
	var b [4]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}
