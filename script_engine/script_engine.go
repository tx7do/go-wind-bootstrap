// Package script_engine provides a bootstrap adapter for the go-scripts script engine.
//
// Import with blank identifier to self-register:
//
//	import _ "github.com/tx7do/go-wind-bootstrap/script_engine"
package script_engine

import (
	"context"
	"fmt"

	scriptEngine "github.com/tx7do/go-scripts"
	"github.com/tx7do/go-scripts/source"

	bootstrap "github.com/tx7do/go-wind-bootstrap"
	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

func init() {
	bootstrap.MustRegisterScriptEngineBuilder("_engine", newBuilder)
}

// newBuilder 根据 Script 配置创建脚本引擎实例。
// 它根据 proto 中的 EngineType 枚举决定创建哪种引擎，
// 并完成 Init、Source 绑定、预加载脚本、热加载等完整初始化流程。
func newBuilder(ctx context.Context, cfg *v1.Script) (any, func(), error) {
	typ := convertEngineType(cfg.GetEngine())
	if typ == "" {
		return nil, nil, fmt.Errorf("script engine: unsupported engine type %v", cfg.GetEngine())
	}

	// 检查 enabled
	if opts := cfg.GetOptions(); opts != nil {
		if !opts.GetEnabled() {
			return nil, func() {}, nil
		}
	}

	// 创建引擎
	eng, err := scriptEngine.NewScriptEngine(typ)
	if err != nil {
		return nil, nil, fmt.Errorf("script engine: create %s failed: %w", typ, err)
	}

	// 初始化
	if err = eng.Init(ctx); err != nil {
		_ = eng.Close()
		return nil, nil, fmt.Errorf("script engine: init %s failed: %w", typ, err)
	}

	// 创建并绑定 Source
	src, err := createSource(cfg.GetSource())
	if err != nil {
		_ = eng.Close()
		return nil, nil, err
	}
	if src != nil {
		eng.SetSource(src)
	}

	// 预加载脚本
	if err = loadScripts(ctx, eng, cfg.GetOptions(), src); err != nil {
		_ = eng.Close()
		return nil, nil, err
	}

	// 热加载
	startHotReload(ctx, eng, cfg.GetOptions())

	cleanup := func() { _ = eng.Close() }
	return eng, cleanup, nil
}

// convertEngineType 将 proto EngineType 枚举转换为 go-scripts Type。
func convertEngineType(t v1.Script_EngineType) scriptEngine.Type {
	switch t {
	case v1.Script_LUA:
		return scriptEngine.LuaType
	case v1.Script_JAVASCRIPT:
		return scriptEngine.JavaScriptType
	case v1.Script_GPYTHON:
		return scriptEngine.GPythonType
	case v1.Script_YAEGI:
		return scriptEngine.YaegiType
	case v1.Script_WAZERO:
		return scriptEngine.WazeroType
	case v1.Script_CEL:
		return scriptEngine.CELType
	case v1.Script_EXPR:
		return scriptEngine.ExprType
	case v1.Script_STARLARK:
		return scriptEngine.StarlarkType
	case v1.Script_TCL:
		return scriptEngine.TclType
	default:
		return ""
	}
}

// createSource 根据配置创建 source.Reader。
func createSource(cfg *v1.Script_Source) (source.Reader, error) {
	if cfg == nil {
		return nil, nil
	}

	var src source.Reader
	var err error

	switch cfg.GetType() {
	case v1.Script_Source_FILE:
		src = source.NewFileSource()

	case v1.Script_Source_MEMORY:
		src = source.NewMemSource()

	case v1.Script_Source_MULTI:
		src, err = createMultiSource(cfg)

	case v1.Script_Source_EMBED:
		return nil, fmt.Errorf("script engine: embed source requires programmatic setup, use SetEmbedFSProvider")

	default:
		// 扩展来源通过工厂注册表
		src, err = createSourceFromFactory(cfg)
	}

	if err != nil {
		return nil, err
	}
	if src == nil {
		return nil, nil
	}

	// 可选缓存层
	if ttlStr := cfg.GetCacheTtl(); ttlStr != "" {
		src, err = wrapWithCache(src, ttlStr)
		if err != nil {
			return nil, err
		}
	}

	return src, nil
}

// loadScripts 加载预加载脚本和入口脚本。
func loadScripts(ctx context.Context, eng any, opts *v1.Script_EngineOptions, src source.Reader) error {
	if opts == nil {
		return nil
	}

	loader := scriptEngine.AsLoader(eng)
	if loader == nil {
		return nil // 轻量引擎不支持文件加载
	}

	paths := opts.GetPaths()

	// 预加载脚本
	for _, script := range opts.GetPreLoadScripts() {
		key := resolveKey(paths, script)
		if src != nil {
			if err := loader.Load(ctx, key); err != nil {
				return fmt.Errorf("script engine: load preload script %q: %w", script, err)
			}
		}
	}

	// 入口脚本
	if entry := opts.GetEntry(); entry != "" {
		key := resolveKey(paths, entry)
		if src != nil {
			if err := loader.Load(ctx, key); err != nil {
				return fmt.Errorf("script engine: load entry script %q: %w", entry, err)
			}
		}
	}

	return nil
}

// startHotReload 启动热加载监听。
func startHotReload(ctx context.Context, eng any, opts *v1.Script_EngineOptions) {
	if opts == nil || !opts.GetHotReload() {
		return
	}

	watcher := scriptEngine.AsWatcher(eng)
	if watcher == nil {
		return
	}

	paths := opts.GetPaths()

	if entry := opts.GetEntry(); entry != "" {
		key := resolveKey(paths, entry)
		_ = watcher.StartWatch(ctx, key)
	}

	for _, script := range opts.GetPreLoadScripts() {
		key := resolveKey(paths, script)
		_ = watcher.StartWatch(ctx, key)
	}
}
