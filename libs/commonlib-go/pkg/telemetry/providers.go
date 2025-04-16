package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
)

type Providers struct {
	TraceProvider  *trace.TracerProvider
	MeterProvider  *metric.MeterProvider
	LoggerProvider *log.LoggerProvider
}

func CreatePropagator() propagation.TextMapPropagator {
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	return propagator
}

func CreateTraceProvider(
	ctx context.Context,
	res *resource.Resource,
	conn *grpc.ClientConn,
) (*trace.TracerProvider, error) {
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}
	batchSpanProcessor := trace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithSpanProcessor(batchSpanProcessor),
		trace.WithBatcher(traceExporter, trace.WithBatchTimeout(5*time.Second)),
	)
	return tracerProvider, nil
}

func CreateMeterProvider(
	ctx context.Context,
	res *resource.Resource,
	conn *grpc.ClientConn,
) (*metric.MeterProvider, error) {
	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create metric exporter: %w", err)
	}
	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExporter, metric.WithInterval(60*time.Second))),
	)
	return meterProvider, nil
}

func CreateLoggerProvider(
	ctx context.Context,
	res *resource.Resource,
	conn *grpc.ClientConn,
) (*log.LoggerProvider, error) {
	logExporter, err := otlploggrpc.New(ctx, otlploggrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create log exporter: %w", err)

	}
	loggerProvider := log.NewLoggerProvider(
		log.WithResource(res),
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
	)
	return loggerProvider, nil
}
