package bootstrap

import (
	"context"
	"sync"

	wind "github.com/tx7do/go-wind"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

// Context holds the application lifecycle state created by [Bootstrap] or [BootstrapWithContext].
//
// It is safe to read fields from multiple goroutines after Bootstrap returns.
type Context struct {
	cfg    *v1.BootstrapConfig
	app    *wind.App
	cancel context.CancelFunc

	brokers  map[string]any
	storages map[string]any

	cleanupOnce sync.Once
	cleanup     func()
}

// newContext creates a Context from the Bootstrap results.
func newContext(cfg *v1.BootstrapConfig, app *wind.App, brokers map[string]any, storages map[string]any, cleanup func(), cancel context.CancelFunc) *Context {
	return &Context{
		cfg:      cfg,
		app:      app,
		brokers:  brokers,
		storages: storages,
		cleanup:  cleanup,
		cancel:   cancel,
	}
}

// Config returns the loaded bootstrap configuration.
func (c *Context) Config() *v1.BootstrapConfig { return c.cfg }

// App returns the underlying [*wind.App].
func (c *Context) App() *wind.App { return c.app }

// Cancel triggers graceful shutdown (idempotent).
func (c *Context) Cancel() {
	if c.cancel != nil {
		c.cancel()
	}
}

// Cleanup releases all resources. It is safe to call multiple times.
func (c *Context) Cleanup() {
	if c == nil {
		return
	}
	c.cleanupOnce.Do(func() {
		if c.cleanup != nil {
			c.cleanup()
		}
	})
}

// Broker returns the broker instance for the given type name (e.g.
// [BrokerTypeKafka], [BrokerTypeRabbitMQ]).
// Returns nil if no broker with that name was configured.
//
// The caller should type-assert the result to the concrete broker type:
//
//	if b, ok := ctx.Broker(bootstrap.BrokerTypeKafka).(*kafka.Broker); ok {
//	    b.Publish(ctx, msg)
//	}
func (c *Context) Broker(name string) any {
	if c == nil || c.brokers == nil {
		return nil
	}
	return c.brokers[name]
}

// Brokers returns all broker instances as a map keyed by type name.
// Returns nil if no broker was configured.
func (c *Context) Brokers() map[string]any {
	if c == nil {
		return nil
	}
	return c.brokers
}

// Storage returns the storage client instance for the given type name (e.g.
// [StorageTypeMinio], [StorageTypeS3]).
// Returns nil if no storage with that name was configured.
//
// The caller should type-assert the result to the concrete storage type:
//
//	s, ok := ctx.Storage(bootstrap.StorageTypeMinio).(*minioPlugin.Storage)
//	if ok { /* use s for PutObject/GetObject */ }
func (c *Context) Storage(name string) any {
	if c == nil || c.storages == nil {
		return nil
	}
	return c.storages[name]
}

// Storages returns all storage instances as a map keyed by type name.
// Returns nil if no storage was configured.
func (c *Context) Storages() map[string]any {
	if c == nil {
		return nil
	}
	return c.storages
}
