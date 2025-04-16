package support

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/config"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/telemetry"
)

func NewTelemetryManager(
	metadata *config.Metadata,
	telemetryCfg *config.TelemetryConfig,
) *TelemetryManager {
	return &TelemetryManager{
		metadata:     metadata,
		telemetryCfg: telemetryCfg,
	}
}

type TelemetryManager struct {
	metadata       *config.Metadata
	telemetryCfg   *config.TelemetryConfig
	TraceProvider  *trace.TracerProvider
	MeterProvider  *metric.MeterProvider
	LoggerProvider *log.LoggerProvider
}

func (m *TelemetryManager) Run() {
	ctx := context.Background()

	res, err := m.createResource(ctx)
	if err != nil {
		panic(err)
	}

	shutdown, providers, err := telemetry.SetupOpenTelemetrySdk(ctx, res, m.telemetryCfg)
	if err != nil {
		panic(err)
	}
	m.TraceProvider = providers.TraceProvider
	m.MeterProvider = providers.MeterProvider
	m.LoggerProvider = providers.LoggerProvider

	defer func() {
		if err := shutdown(ctx); err != nil {
			panic(err)
		}
	}()
}

func (m *TelemetryManager) createResource(ctx context.Context) (*resource.Resource, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(m.metadata.Name),
			semconv.ServiceVersionKey.String(m.metadata.Version),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}
	return res, nil
}
