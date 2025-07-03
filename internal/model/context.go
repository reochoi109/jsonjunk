package model

import (
	"context"
	logger "jsonjunk/pkg/logging"

	"go.uber.org/zap"
)

type contextKey string

const (
	ContextRequestID contextKey = "requestID"
)

func WithContext(ctx context.Context) *zap.Logger {
	if reqID, ok := ctx.Value(ContextRequestID).(string); ok {
		return logger.Log.With(zap.String("requestID", reqID))
	}
	return logger.Log
}
