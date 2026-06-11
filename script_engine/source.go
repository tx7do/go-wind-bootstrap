package script_engine

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/tx7do/go-scripts/source"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

////////////////////////////////////////////////////////////////////////////////
// Source Factory Registry
////////////////////////////////////////////////////////////////////////////////

// SourceFactoryFunc 创建 source.Reader 的工厂函数签名。
type SourceFactoryFunc func(cfg *v1.Script_Source) (source.Reader, error)

var (
	sourceFactoryMu sync.RWMutex
	sourceFactories = make(map[v1.Script_Source_Type]SourceFactoryFunc)
)

// RegisterSourceFactory 注册自定义 Source 创建工厂。
func RegisterSourceFactory(typ v1.Script_Source_Type, f SourceFactoryFunc) error {
	if f == nil {
		return fmt.Errorf("script engine: source factory function is nil")
	}
	sourceFactoryMu.Lock()
	defer sourceFactoryMu.Unlock()
	if _, ok := sourceFactories[typ]; ok {
		return fmt.Errorf("script engine: source factory for %v already registered", typ)
	}
	sourceFactories[typ] = f
	return nil
}

// MustRegisterSourceFactory panics on error。
func MustRegisterSourceFactory(typ v1.Script_Source_Type, f SourceFactoryFunc) {
	if err := RegisterSourceFactory(typ, f); err != nil {
		panic(err)
	}
}

func getSourceFactory(typ v1.Script_Source_Type) (SourceFactoryFunc, bool) {
	sourceFactoryMu.RLock()
	defer sourceFactoryMu.RUnlock()
	f, ok := sourceFactories[typ]
	return f, ok
}

// createSourceFromFactory 通过注册的工厂创建 Source。
func createSourceFromFactory(cfg *v1.Script_Source) (source.Reader, error) {
	f, ok := getSourceFactory(cfg.GetType())
	if !ok {
		return nil, fmt.Errorf(
			"script engine: source type %v not registered — "+
				"import the corresponding source sub-module",
			cfg.GetType(),
		)
	}
	return f(cfg)
}

////////////////////////////////////////////////////////////////////////////////
// Multi Source
////////////////////////////////////////////////////////////////////////////////

func createMultiSource(cfg *v1.Script_Source) (source.Reader, error) {
	paths := cfg.GetPaths()
	if len(paths) == 0 {
		return nil, fmt.Errorf("script engine: multi source requires paths")
	}

	subSources := make([]source.Reader, 0, len(paths))
	for range paths {
		subSources = append(subSources, source.NewFileSource())
	}

	strategy := source.MultiStrategyFallback
	if cfg.GetStrategy() == v1.Script_Source_FIRST_OK {
		strategy = source.MultiStrategyFirstOK
	}

	return source.NewMultiSource(strategy, subSources...)
}

////////////////////////////////////////////////////////////////////////////////
// Cache wrapper
////////////////////////////////////////////////////////////////////////////////

func wrapWithCache(src source.Reader, ttlStr string) (source.Reader, error) {
	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		return nil, fmt.Errorf("script engine: invalid cache_ttl %q: %w", ttlStr, err)
	}
	return source.NewCachedSource(src, source.WithTTL(ttl))
}

////////////////////////////////////////////////////////////////////////////////
// Path resolution
////////////////////////////////////////////////////////////////////////////////

func resolveKey(paths []string, key string) string {
	if key == "" || filepath.IsAbs(key) || len(paths) == 0 {
		return key
	}
	return filepath.Clean(filepath.Join(paths[0], key))
}

// Compile-time assertion.
var _ source.Reader = (*source.FileSource)(nil)

// Export createSourceFromFactory for sub-module access.
var _ = context.Background
