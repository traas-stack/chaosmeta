package log

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"os"
	"strings"
)

const (
	defaultDepth = 3
)

func CtxDebug(ctx context.Context, format string) {
	DefaultLogger.CtxDebugf(ctx, defaultDepth, format)
}

func CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	DefaultLogger.CtxDebugf(ctx, defaultDepth, format, v...)
}

func CtxInfo(ctx context.Context, format string) {
	DefaultLogger.CtxInfof(ctx, defaultDepth, format)
}

func CtxInfof(ctx context.Context, format string, v ...interface{}) {
	DefaultLogger.CtxInfof(ctx, defaultDepth, format, v...)
}

func CtxWarn(ctx context.Context, format string) {
	DefaultLogger.CtxWarningf(ctx, defaultDepth, format)
}

func CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	DefaultLogger.CtxWarningf(ctx, defaultDepth, format, v...)
}

func CtxError(ctx context.Context, format string) {
	DefaultLogger.CtxErrorf(ctx, defaultDepth, format)
}

func CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	DefaultLogger.CtxErrorf(ctx, defaultDepth, format, v...)
}

func CtxFatal(ctx context.Context, format string) {
	DefaultLogger.CtxFatalf(ctx, defaultDepth, format)
}

func CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	DefaultLogger.CtxFatalf(ctx, defaultDepth, format, v...)
}

func CtxPanic(ctx context.Context, format string) {
	DefaultLogger.CtxPanicf(ctx, defaultDepth, format)
}

func CtxPanicf(ctx context.Context, format string, v ...interface{}) {
	DefaultLogger.CtxPanicf(ctx, defaultDepth, format, v...)
}

func (l *LoggerStruct) CtxFatal(ctx context.Context, depth int, format string, v ...interface{}) {
	l.ctxlogf(ctx, depth+1, "fatal", format, v...)
}

func (l *LoggerStruct) CtxFatalf(ctx context.Context, depth int, format string, v ...interface{}) {
	l.ctxlogf(ctx, depth+1, "fatal", format, v...)
}

func (l *LoggerStruct) CtxPanic(ctx context.Context, depth int, format string, v ...interface{}) {
	l.ctxlogf(ctx, depth+1, "panic", format, v...)
	os.Exit(-1)
}

func (l *LoggerStruct) CtxPanicf(ctx context.Context, depth int, format string, v ...interface{}) {
	l.ctxlogf(ctx, depth+1, "panic", format, v...)
	os.Exit(-1)
}

func (l *LoggerStruct) CtxError(ctx context.Context, depth int, format string, v ...interface{}) {
	l.ctxlogf(ctx, depth+1, "error", format, v...)
}

func (l *LoggerStruct) CtxErrorf(ctx context.Context, depth int, format string, v ...interface{}) {
	l.ctxlogf(ctx, depth+1, "error", format, v...)
}

func (l *LoggerStruct) CtxWarning(ctx context.Context, depth int, format string, v ...interface{}) {
	l.ctxlogf(ctx, depth+1, "warn", format, v...)
}

func (l *LoggerStruct) CtxWarningf(ctx context.Context, depth int, format string, v ...interface{}) {
	l.ctxlogf(ctx, depth+1, "warn", format, v...)
}

func (l *LoggerStruct) CtxInfo(ctx context.Context, depth int, format string, v ...interface{}) {
	l.ctxlogf(ctx, depth+1, "info", format, v...)
}

func (l *LoggerStruct) CtxInfof(ctx context.Context, depth int, format string, v ...interface{}) {
	l.ctxlogf(ctx, depth+1, "info", format, v...)
}

func (l *LoggerStruct) CtxDebug(ctx context.Context, depth int, format string, v ...interface{}) {
	l.ctxlogf(ctx, depth+1, "debug", format, v...)
}

func (l *LoggerStruct) CtxDebugf(ctx context.Context, depth int, format string, v ...interface{}) {
	l.ctxlogf(ctx, depth+1, "debug", format, v...)
}

func (l *LoggerStruct) ctxlogf(ctx context.Context, depth int, t string, format string, v ...interface{}) {
	if ctx == nil {
		l.outPutByLevel(depth, t, fmt.Sprintf(format, v...))
		return
	}
	traceVal := ctx.Value(TraceIdKey)
	traceInf, ok := traceVal.(trace)
	if !ok {
		l.outPutByLevel(depth, t, fmt.Sprintf(format, v...))
		return
	}

	traceStr := traceInf.marshal()
	l.outPutByLevel(depth, t, fmt.Sprintf("[%s] %s", traceStr, fmt.Sprintf(format, v...)))
}

func (l *LoggerStruct) outPutByLevel(depth int, level, msg string) {
	switch strings.ToLower(level) {
	case "fatal":
		l.globalLogger.WithOptions(zap.AddCallerSkip(depth)).Fatal(msg)
	case "panic":
		l.globalLogger.WithOptions(zap.AddCallerSkip(depth)).Panic(msg)
	case "error":
		l.globalLogger.WithOptions(zap.AddCallerSkip(depth)).Error(msg)
	case "warn":
		l.globalLogger.WithOptions(zap.AddCallerSkip(depth)).Warn(msg)
	case "info":
		l.globalLogger.WithOptions(zap.AddCallerSkip(depth)).Info(msg)
	case "debug":
		l.globalLogger.WithOptions(zap.AddCallerSkip(depth)).Debug(msg)
	default:
		l.globalLogger.WithOptions(zap.AddCallerSkip(depth)).Debug(msg)
	}
}

func getTraceId(ctx context.Context) string {
	if ctx.Value(TraceIdKey) == nil {
		return "[trace: system]"
	}

	return fmt.Sprintf("[trace: %s]", ctx.Value(TraceIdKey).(string))
}
