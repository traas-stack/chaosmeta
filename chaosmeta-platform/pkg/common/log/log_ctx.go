package log

import (
	"context"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
)

const spaceStr = " "

func DebugCtx(ctx context.Context, v ...interface{}) {
	logs.Debug(fmt.Sprintf("%s %v", getTraceId(ctx)+spaceStr, v))
}

func DebugCtxf(ctx context.Context, f string, v ...interface{}) {
	logs.Debug(getTraceId(ctx) + spaceStr + fmt.Sprintf(f, v...))
}

func InfoCtx(ctx context.Context, v ...interface{}) {
	logs.Info(fmt.Sprintf("%s %v", getTraceId(ctx)+spaceStr, v))
}

func InfoCtxf(ctx context.Context, f string, v ...interface{}) {
	logs.Info(getTraceId(ctx) + spaceStr + fmt.Sprintf(f, v...))
}

func WarnCtx(ctx context.Context, v ...interface{}) {
	logs.Warn(fmt.Sprintf("%s %v", getTraceId(ctx)+spaceStr, v))
}

func WarnCtxf(ctx context.Context, f string, v ...interface{}) {
	logs.Warn(getTraceId(ctx) + spaceStr + fmt.Sprintf(f, v...))
}

func ErrorCtx(ctx context.Context, v ...interface{}) {
	logs.Error(fmt.Sprintf("%s %v", getTraceId(ctx)+spaceStr, v))
}

func ErrorCtxf(ctx context.Context, f string, v ...interface{}) {
	logs.Error(getTraceId(ctx) + spaceStr + fmt.Sprintf(f, v...))
}

func ErrorfCtx(ctx context.Context, msg string) {
	logs.Error(fmt.Sprintf("[trace: %s] %s", getTraceId(ctx), msg))
}

func getTraceId(ctx context.Context) string {
	if ctx.Value(TraceIdKey) == nil {
		return "[trace: system]"
	}

	return fmt.Sprintf("[trace: %s]", ctx.Value(TraceIdKey).(string))
}
