package script_engine

import (
	"context"

	scriptEngine "github.com/tx7do/go-scripts"
	windLog "github.com/tx7do/go-wind/log"
)

// Compile-time assertion: WindLogger implements scriptEngine.Logger.
var _ scriptEngine.Logger = (*WindLogger)(nil)

// WindLogger 将 go-wind log.Logger 适配为 go-scripts Logger 接口。
//
// 使得 go-scripts 内部日志（引擎初始化、脚本执行、热更新等）能够
// 无缝接入 go-wind 统一日志体系。
type WindLogger struct {
	Logger windLog.Logger
}

// Debug 以 DEBUG 级别记录日志。
func (l *WindLogger) Debug(ctx context.Context, msg string, args ...any) {
	if l.Logger == nil {
		return
	}
	l.Logger.Debug(ctx, msg, args...)
}

// Info 以 INFO 级别记录日志。
func (l *WindLogger) Info(ctx context.Context, msg string, args ...any) {
	if l.Logger == nil {
		return
	}
	l.Logger.Info(ctx, msg, args...)
}

// Warn 以 WARN 级别记录日志。
func (l *WindLogger) Warn(ctx context.Context, msg string, args ...any) {
	if l.Logger == nil {
		return
	}
	l.Logger.Warn(ctx, msg, args...)
}

// Error 以 ERROR 级别记录日志。
func (l *WindLogger) Error(ctx context.Context, msg string, args ...any) {
	if l.Logger == nil {
		return
	}
	l.Logger.Error(ctx, msg, args...)
}

// With 返回附加了指定键值对的新 WindLogger 实例。
func (l *WindLogger) With(args ...any) scriptEngine.Logger {
	if l.Logger == nil {
		return l
	}
	return &WindLogger{Logger: l.Logger.With(args...)}
}

// SetWindLogger 设置 go-scripts 全局日志为 go-wind Logger。
// 传入 nil 可恢复为默认的静默日志。
func SetWindLogger(windLogger windLog.Logger) {
	if windLogger == nil {
		scriptEngine.SetLogger(nil)
		return
	}
	scriptEngine.SetLogger(&WindLogger{Logger: windLogger})
}
