package logger

import (
	"context"
	"time"

	"github.com/ritchieridanko/erteku/services/auth/internal/utils"
	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.Logger
}

func NewLogger(l *zap.Logger) *Logger {
	return &Logger{logger: l}
}

func (l *Logger) Log(message string, args ...any) {
	l.logger.Sugar().Infof(message, args...)
}

func (l *Logger) Info(ctx context.Context, message string, fields ...Field) {
	l.logger.Info(message, l.toFields(ctx, fields...)...)
}

func (l *Logger) Warn(ctx context.Context, message string, fields ...Field) {
	l.logger.Warn(message, l.toFields(ctx, fields...)...)
}

func (l *Logger) Error(ctx context.Context, message string, fields ...Field) {
	l.logger.Error(message, l.toFields(ctx, fields...)...)
}

func (l *Logger) toFields(ctx context.Context, fields ...Field) []zap.Field {
	zf := make([]zap.Field, 0, len(fields)+3)
	zf = append(zf, zap.Time("timestamp", time.Now().UTC()))

	if requestID := utils.CtxRequestID(ctx); requestID != "" {
		zf = append(zf, zap.String("request_id", requestID))
	}
	if traceID := utils.CtxTraceID(ctx); traceID != "" {
		zf = append(zf, zap.String("trace_id", traceID))
	}

	for _, field := range fields {
		switch v := field.value.(type) {
		case error:
			zf = append(zf, zap.Error(v))
		default:
			zf = append(zf, zap.Any(field.key, field.value))
		}
	}
	return zf
}
