package bootstrap

import (
	"context"
	"sync"

	wind "github.com/tx7do/go-wind"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

// Context holds the application lifecycle state created by [Bootstrap] or [RunApp].
//
// It is safe to read fields from multiple goroutines after Bootstrap returns.
type Context struct {
	cfg    *v1.BootstrapConfig
	app    *wind.App
	cancel context.CancelFunc

	cleanupOnce sync.Once
	cleanup     func()
}

// newContext creates a Context from the Bootstrap results.
func newContext(cfg *v1.BootstrapConfig, app *wind.App, cleanup func(), cancel context.CancelFunc) *Context {
	return &Context{
		cfg:     cfg,
		app:     app,
		cleanup: cleanup,
		cancel:  cancel,
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
