// Package bootstrap provides builder registration for all plugin domains.
//
// Each plugin domain (server, config, registry, log, tracer, metrics, broker)
// maintains a map of string-keyed builder functions. The key is a lowercase
// type string (e.g. "consul", "zap") matching the JSON config value.
// Built-in builders are registered via init() in provider packages.
// Users can register custom builders to extend the framework.
package bootstrap

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/tx7do/go-wind/log"
	"github.com/tx7do/go-wind/transport"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

// ServerBuilder builds a [transport.Server] from a [Server] config.
type ServerBuilder func(cfg *v1.Server) (transport.Server, error)

// LogBuilder builds a [log.Logger] from a [Logger] config.
type LogBuilder func(cfg *v1.Logger) (log.Logger, func(), error)

// RegistryAction performs registration/deregistration lifecycle.
//
// The cfg parameter carries the full registry configuration so that
// the action can extract its own sub-configuration.
// Deprecated: Use RegisterRegistryAction with the new signature instead.
type RegistryAction func(ctx context.Context, appCfg *v1.App, endpoints []string, cfg *v1.Registry) (func(), error)

// ConfigAction performs config source loading/watching.
type ConfigAction func(ctx context.Context, cfg *v1.Config) (func(), error)

// TracerBuilder builds a tracer provider.
type TracerBuilder func(cfg *v1.Tracer) (interface{}, func(), error)

// MetricsBuilder builds a metrics backend.
type MetricsBuilder func(cfg *v1.Metrics) (func(), error)

// BrokerBuilder builds a broker instance and returns it along with an optional
// cleanup function. The returned instance (any) is the concrete broker object
// that callers can use for Publish/Subscribe operations.
//
// The type key (e.g. "kafka", "rabbitmq") is used to look up the instance
// via [Context.Broker] after bootstrap.
type BrokerBuilder func(ctx context.Context, cfg *v1.Broker) (any, func(), error)

// StorageBuilder builds a storage client instance and returns it along with an
// optional cleanup function. The returned instance (any) is the concrete
// storage client that callers can use for object operations.
//
// The type key (e.g. "minio", "s3") is used to look up the instance
// via [Context.Storage] after bootstrap.
type StorageBuilder func(ctx context.Context, cfg *v1.Storage) (any, func(), error)

// ---- Global registries (string keyed) ----

var (
	mu              sync.RWMutex
	serverBuilders  = map[string]ServerBuilder{}
	logBuilders     = map[string]LogBuilder{}
	registryActions = map[string]RegistryAction{}
	configActions   = map[string]ConfigAction{}
	tracerBuilders  = map[string]TracerBuilder{}
	metricsBuilders = map[string]MetricsBuilder{}
	brokerBuilders  = map[string]BrokerBuilder{}
	storageBuilders = map[string]StorageBuilder{}
)

// ---- Register functions ----

// RegisterServerBuilder registers a server builder for the given type string.
func RegisterServerBuilder(typ string, b ServerBuilder) error {
	if typ == "" {
		return fmt.Errorf("bootstrap: type is empty")
	}
	if b == nil {
		return fmt.Errorf("bootstrap: factory is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := serverBuilders[typ]; ok {
		return fmt.Errorf("bootstrap: server builder %q already registered", typ)
	}
	serverBuilders[typ] = b
	return nil
}

// MustRegisterServerBuilder panics on error. Intended for init().
func MustRegisterServerBuilder(typ string, b ServerBuilder) {
	if err := RegisterServerBuilder(typ, b); err != nil {
		panic(err)
	}
}

// RegisterLogBuilder registers a log builder for the given type string.
func RegisterLogBuilder(typ string, b LogBuilder) error {
	if typ == "" {
		return fmt.Errorf("bootstrap: type is empty")
	}
	if b == nil {
		return fmt.Errorf("bootstrap: factory is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := logBuilders[typ]; ok {
		return fmt.Errorf("bootstrap: log builder %q already registered", typ)
	}
	logBuilders[typ] = b
	return nil
}

// MustRegisterLogBuilder panics on error.
func MustRegisterLogBuilder(typ string, b LogBuilder) {
	if err := RegisterLogBuilder(typ, b); err != nil {
		panic(err)
	}
}

// RegisterRegistryAction registers a registry action for the given type string.
func RegisterRegistryAction(typ string, a RegistryAction) error {
	if typ == "" {
		return fmt.Errorf("bootstrap: type is empty")
	}
	if a == nil {
		return fmt.Errorf("bootstrap: factory is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := registryActions[typ]; ok {
		return fmt.Errorf("bootstrap: registry action %q already registered", typ)
	}
	registryActions[typ] = a
	return nil
}

// MustRegisterRegistryAction panics on error.
func MustRegisterRegistryAction(typ string, a RegistryAction) {
	if err := RegisterRegistryAction(typ, a); err != nil {
		panic(err)
	}
}

// RegisterConfigAction registers a config action for the given type string.
func RegisterConfigAction(typ string, a ConfigAction) error {
	if typ == "" {
		return fmt.Errorf("bootstrap: type is empty")
	}
	if a == nil {
		return fmt.Errorf("bootstrap: factory is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := configActions[typ]; ok {
		return fmt.Errorf("bootstrap: config action %q already registered", typ)
	}
	configActions[typ] = a
	return nil
}

// MustRegisterConfigAction panics on error.
func MustRegisterConfigAction(typ string, a ConfigAction) {
	if err := RegisterConfigAction(typ, a); err != nil {
		panic(err)
	}
}

// RegisterTracerBuilder registers a tracer builder for the given type string.
func RegisterTracerBuilder(typ string, b TracerBuilder) error {
	if typ == "" {
		return fmt.Errorf("bootstrap: type is empty")
	}
	if b == nil {
		return fmt.Errorf("bootstrap: factory is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := tracerBuilders[typ]; ok {
		return fmt.Errorf("bootstrap: tracer builder %q already registered", typ)
	}
	tracerBuilders[typ] = b
	return nil
}

// MustRegisterTracerBuilder panics on error.
func MustRegisterTracerBuilder(typ string, b TracerBuilder) {
	if err := RegisterTracerBuilder(typ, b); err != nil {
		panic(err)
	}
}

// RegisterMetricsBuilder registers a metrics builder for the given type string.
func RegisterMetricsBuilder(typ string, b MetricsBuilder) error {
	if typ == "" {
		return fmt.Errorf("bootstrap: type is empty")
	}
	if b == nil {
		return fmt.Errorf("bootstrap: factory is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := metricsBuilders[typ]; ok {
		return fmt.Errorf("bootstrap: metrics builder %q already registered", typ)
	}
	metricsBuilders[typ] = b
	return nil
}

// MustRegisterMetricsBuilder panics on error.
func MustRegisterMetricsBuilder(typ string, b MetricsBuilder) {
	if err := RegisterMetricsBuilder(typ, b); err != nil {
		panic(err)
	}
}

// RegisterBrokerBuilder registers a broker builder for the given type string.
func RegisterBrokerBuilder(typ string, b BrokerBuilder) error {
	if typ == "" {
		return fmt.Errorf("bootstrap: type is empty")
	}
	if b == nil {
		return fmt.Errorf("bootstrap: factory is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := brokerBuilders[typ]; ok {
		return fmt.Errorf("bootstrap: broker builder %q already registered", typ)
	}
	brokerBuilders[typ] = b
	return nil
}

// MustRegisterBrokerBuilder panics on error.
func MustRegisterBrokerBuilder(typ string, b BrokerBuilder) {
	if err := RegisterBrokerBuilder(typ, b); err != nil {
		panic(err)
	}
}

// RegisterStorageBuilder registers a storage builder for the given type string.
func RegisterStorageBuilder(typ string, b StorageBuilder) error {
	if typ == "" {
		return fmt.Errorf("bootstrap: type is empty")
	}
	if b == nil {
		return fmt.Errorf("bootstrap: factory is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := storageBuilders[typ]; ok {
		return fmt.Errorf("bootstrap: storage builder %q already registered", typ)
	}
	storageBuilders[typ] = b
	return nil
}

// MustRegisterStorageBuilder panics on error.
func MustRegisterStorageBuilder(typ string, b StorageBuilder) {
	if err := RegisterStorageBuilder(typ, b); err != nil {
		panic(err)
	}
}

// ---- Lookup helpers ----

func getServerBuilder(typ string) (ServerBuilder, error) {
	mu.RLock()
	defer mu.RUnlock()
	b, ok := serverBuilders[typ]
	if !ok {
		return nil, fmt.Errorf("bootstrap: no server builder registered for %q", typ)
	}
	return b, nil
}

func getLogBuilder(typ string) (LogBuilder, error) {
	mu.RLock()
	defer mu.RUnlock()
	b, ok := logBuilders[typ]
	if !ok {
		return nil, fmt.Errorf("bootstrap: no log builder registered for %q", typ)
	}
	return b, nil
}

func getRegistryAction(typ string) (RegistryAction, error) {
	mu.RLock()
	defer mu.RUnlock()
	a, ok := registryActions[typ]
	if !ok {
		return nil, fmt.Errorf("bootstrap: no registry action registered for %q", typ)
	}
	return a, nil
}

func getConfigAction(typ string) (ConfigAction, error) {
	mu.RLock()
	defer mu.RUnlock()
	a, ok := configActions[typ]
	if !ok {
		return nil, fmt.Errorf("bootstrap: no config action registered for %q", typ)
	}
	return a, nil
}

func getTracerBuilder(typ string) (TracerBuilder, error) {
	mu.RLock()
	defer mu.RUnlock()
	b, ok := tracerBuilders[typ]
	if !ok {
		return nil, fmt.Errorf("bootstrap: no tracer builder registered for %q", typ)
	}
	return b, nil
}

func getMetricsBuilder(typ string) (MetricsBuilder, error) {
	mu.RLock()
	defer mu.RUnlock()
	b, ok := metricsBuilders[typ]
	if !ok {
		return nil, fmt.Errorf("bootstrap: no metrics builder registered for %q", typ)
	}
	return b, nil
}

func getBrokerBuilder(typ string) (BrokerBuilder, error) {
	mu.RLock()
	defer mu.RUnlock()
	b, ok := brokerBuilders[typ]
	if !ok {
		return nil, fmt.Errorf("bootstrap: no broker builder registered for %q", typ)
	}
	return b, nil
}

func getStorageBuilder(typ string) (StorageBuilder, error) {
	mu.RLock()
	defer mu.RUnlock()
	b, ok := storageBuilders[typ]
	if !ok {
		return nil, fmt.Errorf("bootstrap: no storage builder registered for %q", typ)
	}
	return b, nil
}

// ---- Diagnostic helpers ----

// ListServerBuilders returns all registered server type names.
func ListServerBuilders() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(serverBuilders))
	for k := range serverBuilders {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// ListLogBuilders returns all registered log type names.
func ListLogBuilders() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(logBuilders))
	for k := range logBuilders {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// ListRegistryActions returns all registered registry type names.
func ListRegistryActions() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(registryActions))
	for k := range registryActions {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// ListConfigActions returns all registered config type names.
func ListConfigActions() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(configActions))
	for k := range configActions {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// ListTracerBuilders returns all registered tracer type names.
func ListTracerBuilders() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(tracerBuilders))
	for k := range tracerBuilders {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// ListMetricsBuilders returns all registered metrics type names.
func ListMetricsBuilders() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(metricsBuilders))
	for k := range metricsBuilders {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// ListBrokerBuilders returns all registered broker type names.
func ListBrokerBuilders() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(brokerBuilders))
	for k := range brokerBuilders {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// ListStorageBuilders returns all registered storage type names.
func ListStorageBuilders() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(storageBuilders))
	for k := range storageBuilders {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}
