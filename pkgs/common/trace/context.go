package trace

import (
	"context"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type contextKey string

const correlationIDKey contextKey = "correlation_id"

// GetCorrelationID extracts correlation ID from context.
func GetCorrelationID(ctx context.Context) string {
	if id, ok := ctx.Value(correlationIDKey).(string); ok {
		return id
	}
	return ""
}

// WithCorrelationID adds correlation ID to context.
func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, correlationIDKey, id)
}

// ContextFromNATSMsg builds a context with trace and correlation ID from NATS headers.
func ContextFromNATSMsg(msg *nats.Msg) context.Context {
	ctx := context.Background()
	if msg == nil {
		return ctx
	}
	if msg.Header != nil {
		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(msg.Header))
		if corrID := msg.Header.Get("X-Correlation-ID"); corrID != "" {
			ctx = WithCorrelationID(ctx, corrID)
		}
	}
	return ctx
}

// InjectContextHeaders writes trace and correlation ID into NATS message headers.
func InjectContextHeaders(ctx context.Context, header nats.Header) {
	if header == nil {
		return
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(header))
	if corrID := GetCorrelationID(ctx); corrID != "" {
		header.Set("X-Correlation-ID", corrID)
	}
}