package application

import (
	"context"
	"errors"

	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

type openTelemetry struct {
	tp *trace.TracerProvider
	mp *metric.MeterProvider
}

func (op *openTelemetry) init(ctx context.Context) error {
	if err := op.initOT(ctx); err != nil {
		return err
	}

	op.initGlobal()
	return nil
}

func (op *openTelemetry) initOT(ctx context.Context) error {
	// Create trace exporter using environment variables
	spanExporter, err := autoexport.NewSpanExporter(ctx)
	if err != nil {
		return err
	}

	// Create metric reader using environment variables
	metricReader, err := autoexport.NewMetricReader(ctx)
	if err != nil {
		return err
	}

	// Create trace provider with the exporter
	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(spanExporter),
	)

	// Create meter provider with the reader
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metricReader),
	)

	op.tp = tracerProvider
	op.mp = meterProvider
	return nil
}

func (ot *openTelemetry) initGlobal() {
	otel.SetTracerProvider(ot.tp)
	otel.SetMeterProvider(ot.mp)
}

func (ot *openTelemetry) shutdown(ctx context.Context) error {
	return errors.Join(
		ot.tp.Shutdown(ctx),
		ot.mp.Shutdown(ctx),
	)
}
