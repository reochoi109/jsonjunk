package model

import (
	"context"
	logger "jsonjunk/pkg/logging"

	"go.uber.org/zap"
)

type contextKey string

const (
	ContextTraceID contextKey = "trace_id"
)

func WithContext(ctx context.Context) *zap.Logger {
	if reqID, ok := ctx.Value(ContextTraceID).(string); ok {
		return logger.Log.With(zap.String("trace_id", reqID))
	}
	return logger.Log
}
