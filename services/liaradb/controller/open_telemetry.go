package controller

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	actionAttribute        = "request.action"
	aggregateIDAttribute   = "request.aggregate_id"
	correlationIDAttribute = "request.correlation_id"
	tenantIDAttribute      = "request.tenant_id"
)

func traceAppend(ctx context.Context, tenantID string, correlationID string, aggregateIDs []string) {
	span := trace.SpanFromContext(ctx)

	span.SetAttributes(
		attribute.String(actionAttribute, "append"),
		attribute.String(tenantIDAttribute, tenantID),
		attribute.String(correlationIDAttribute, correlationID),
		attribute.StringSlice(aggregateIDAttribute, aggregateIDs))
}
