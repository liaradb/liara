package controller

import (
	"context"

	"github.com/liaradb/liaradb/domain/value"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	correlationIDAttribute = "correlation_id"
)

type CorrelationIDer interface {
	CorrelationID() value.CorrelationID
}

func traceCorrelationID(ctx context.Context, r CorrelationIDer) {
	cid := r.CorrelationID().String()
	if cid == "" {
		span := trace.SpanFromContext(ctx)
		span.SetAttributes(attribute.String(correlationIDAttribute, cid))
	}
}
