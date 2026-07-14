package common

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/pafthang/pocketagent/pkgs/common/trace"
)

func GetCorrelationID(ctx context.Context) string { return trace.GetCorrelationID(ctx) }
func WithCorrelationID(ctx context.Context, id string) context.Context {
	return trace.WithCorrelationID(ctx, id)
}
func ContextFromNATSMsg(msg *nats.Msg) context.Context { return trace.ContextFromNATSMsg(msg) }
func InjectContextHeaders(ctx context.Context, header nats.Header) {
	trace.InjectContextHeaders(ctx, header)
}
func RootCorrelationID(corrID string) string          { return trace.RootCorrelationID(corrID) }
func SubtaskIndex(parentCorrID, subCorrID string) int { return trace.SubtaskIndex(parentCorrID, subCorrID) }