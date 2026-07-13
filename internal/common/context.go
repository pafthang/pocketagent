package common

import (
	"context"
)

// Context keys
type contextKey string

const (
	CorrelationIDKey contextKey = "correlation_id"
	TraceIDKey       contextKey = "trace_id"
)

// GetCorrelationID extracts correlation ID from context
func GetCorrelationID(ctx context.Context) string {
	if id, ok := ctx.Value(CorrelationIDKey).(string); ok {
		return id
	}
	return ""
}

// WithCorrelationID adds correlation ID to context
func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, CorrelationIDKey, id)
}
