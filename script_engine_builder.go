package bootstrap

import (
	"context"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
)

// resolveScriptEngine 解析 Script 配置，查找对应的 ScriptEngineBuilder 并调用。
// 返回按引擎类型名索引的实例映射和统一的 cleanup 函数。
func resolveScriptEngine(ctx context.Context, cfg *v1.Script) (map[string]any, func(), error) {
	if cfg == nil {
		return nil, nil, nil
	}

	// 从 proto EngineType 枚举映射到 string key。
	typ := convertScriptEngineType(cfg.GetEngine())
	if typ == "" {
		return nil, nil, nil
	}

	// 检查 enabled 标志
	if opts := cfg.GetOptions(); opts != nil {
		if !opts.GetEnabled() {
			return nil, nil, nil
		}
	}

	b, err := getScriptEngineBuilder(typ)
	if err != nil {
		return nil, nil, err
	}

	inst, cleanup, err := b(ctx, cfg)
	if err != nil {
		return nil, nil, err
	}

	instances := map[string]any{typ: inst}
	return instances, cleanup, nil
}

// convertScriptEngineType 将 proto EngineType 枚举转换为注册表 key 字符串。
func convertScriptEngineType(t v1.Script_EngineType) string {
	switch t {
	case v1.Script_LUA:
		return ScriptEngineLua
	case v1.Script_JAVASCRIPT:
		return ScriptEngineJavaScript
	case v1.Script_GPYTHON:
		return ScriptEngineGPython
	case v1.Script_YAEGI:
		return ScriptEngineYaegi
	case v1.Script_WAZERO:
		return ScriptEngineWazero
	case v1.Script_CEL:
		return ScriptEngineCEL
	case v1.Script_EXPR:
		return ScriptEngineExpr
	case v1.Script_STARLARK:
		return ScriptEngineStarlark
	case v1.Script_TCL:
		return ScriptEngineTcl
	default:
		return ""
	}
}
