package log

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	TraceIdKey = "TraceId"
	sepSig     = ";"
)

type trace struct {
	TraceId string `json:"trace_id,omitempty"`
	Tips    string `json:"tips,omitempty"`
}

func (t trace) marshal() string {
	bts, err := json.Marshal(t)
	if err != nil {
		return ""
	}

	return string(bts)
}

func TraceCtx(ctx context.Context, traceId string) context.Context {
	return context.WithValue(ctx, TraceIdKey, trace{
		TraceId: traceId,
	})
}

func AppendTipsCtx(ctx context.Context, tips string) context.Context {
	rawVal := ctx.Value(TraceIdKey)
	rawTrace, ok := rawVal.(trace)
	if !ok {
		return context.WithValue(ctx, TraceIdKey, trace{
			Tips: tips})
	}

	tips = fmt.Sprintf("%s%s%s", rawTrace.Tips, tips, sepSig)

	return context.WithValue(ctx, TraceIdKey, trace{
		TraceId: rawTrace.TraceId,
		Tips:    tips,
	})
}

func TraceAndTipCtx(ctx context.Context, traceId, tips string) context.Context {
	ctx = TraceCtx(ctx, traceId)
	ctx = AppendTipsCtx(ctx, tips)

	return ctx
}
