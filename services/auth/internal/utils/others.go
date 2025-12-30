package utils

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/ritchieridanko/erteku/services/auth/internal/constants"
	"go.opentelemetry.io/otel/trace"
)

// Get request id from context
func CtxRequestID(ctx context.Context) string {
	return ctx.Value(constants.CtxKeyRequestID).(string)
}

// Get request meta (user agent and IP address) from context
func CtxRequestMeta(ctx context.Context) (userAgent, ipAddress string) {
	return ctx.Value(constants.CtxKeyUserAgent).(string), ctx.Value(constants.CtxKeyIPAddress).(string)
}

// Get trace id from context
func CtxTraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().HasTraceID() {
		return span.SpanContext().TraceID().String()
	}

	return ""
}

// Create a new random UUID
func GenerateUUID() uuid.UUID {
	return uuid.New()
}

// Remove string of leading and trailing whitespaces
// and set to all lowercase
func NormalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
