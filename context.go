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

	brokers       map[string]any
	storages      map[string]any
	aiClients     map[string]any
	workflows     map[string]any
	caches        map[string]any
	scriptEngines map[string]any
	databases     map[string]any

	cleanupOnce sync.Once
	cleanup     func()
}

// newContext creates a Context from the Bootstrap results.
func newContext(cfg *v1.BootstrapConfig, app *wind.App, brokers map[string]any, storages map[string]any, aiClients map[string]any, workflows map[string]any, caches map[string]any, scriptEngines map[string]any, databases map[string]any, cleanup func(), cancel context.CancelFunc) *Context {
	return &Context{
		cfg:           cfg,
		app:           app,
		brokers:       brokers,
		storages:      storages,
		aiClients:     aiClients,
		workflows:     workflows,
		caches:        caches,
		scriptEngines: scriptEngines,
		databases:     databases,
		cleanup:       cleanup,
		cancel:        cancel,
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

// Ai returns the AI model client instance for the given type name (e.g.
// [AiTypeOpenAI], [AiTypeLangChainGo], [AiTypeEino]).
// Returns nil if no AI client with that name was configured.
//
// The caller should type-assert the result to the concrete client type:
//
//	client, ok := ctx.Ai(bootstrap.AiTypeOpenAI).(*openai.Client)
//	if ok { /* use client for chat completions */ }
func (c *Context) Ai(name string) any {
	if c == nil || c.aiClients == nil {
		return nil
	}
	return c.aiClients[name]
}

// Ais returns all AI client instances as a map keyed by type name.
// Returns nil if no AI client was configured.
func (c *Context) Ais() map[string]any {
	if c == nil {
		return nil
	}
	return c.aiClients
}

// Workflow returns the workflow engine client instance for the given type name
// (e.g. [WorkflowTypeTemporal], [WorkflowTypeArgo]).
// Returns nil if no workflow client with that name was configured.
//
// The caller should type-assert the result to the concrete client type:
//
//	wc, ok := ctx.Workflow(bootstrap.WorkflowTypeTemporal).(*temporal.WorkflowClient)
//	if ok { /* use wc for workflow operations */ }
func (c *Context) Workflow(name string) any {
	if c == nil || c.workflows == nil {
		return nil
	}
	return c.workflows[name]
}

// Workflows returns all workflow client instances as a map keyed by type name.
// Returns nil if no workflow client was configured.
func (c *Context) Workflows() map[string]any {
	if c == nil {
		return nil
	}
	return c.workflows
}

// Cache returns the cache instance for the given type name (e.g.
// [CacheTypeLocal], [CacheTypeRedis]).
// Returns nil if no cache with that name was configured.
//
// The caller should type-assert the result to the concrete cache type:
//
//	c, ok := ctx.Cache(bootstrap.CacheTypeLocal).(*local.Cache)
//	if ok { /* use c for Get/Set operations */ }
func (c *Context) Cache(name string) any {
	if c == nil || c.caches == nil {
		return nil
	}
	return c.caches[name]
}

// Caches returns all cache instances as a map keyed by type name.
// Returns nil if no cache was configured.
func (c *Context) Caches() map[string]any {
	if c == nil {
		return nil
	}
	return c.caches
}

// ScriptEngine returns the script engine instance for the given type name (e.g.
// [ScriptEngineLua], [ScriptEngineJavaScript]).
// Returns nil if no script engine with that name was configured.
//
// The caller should type-assert the result to the concrete engine type:
//
//	eng, ok := ctx.ScriptEngine(bootstrap.ScriptEngineLua).(scriptEngine.Engine)
//	if ok { /* use eng for Eval/Execute */ }
func (c *Context) ScriptEngine(name string) any {
	if c == nil || c.scriptEngines == nil {
		return nil
	}
	return c.scriptEngines[name]
}

// ScriptEngines returns all script engine instances as a map keyed by type name.
// Returns nil if no script engine was configured.
func (c *Context) ScriptEngines() map[string]any {
	if c == nil {
		return nil
	}
	return c.scriptEngines
}

// Database returns the database client instance for the given type name (e.g.
// [DatabaseTypeGorm], [DatabaseTypeMongodb]).
// Returns nil if no database client with that name was configured.
func (c *Context) Database(name string) any {
	if c == nil || c.databases == nil {
		return nil
	}
	return c.databases[name]
}

// Databases returns all database client instances as a map keyed by type name.
// Returns nil if no database was configured.
func (c *Context) Databases() map[string]any {
	if c == nil {
		return nil
	}
	return c.databases
}
