package support

import (
	"context"
	"fmt"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
	"go.uber.org/fx"

	"github.com/codejsha/shared-library-go/pkg/config"
	"github.com/codejsha/shared-library-go/pkg/telemetry"
)

type TelemetryManager struct {
	metadata       *config.Metadata
	telemetryCfg   *config.TelemetryConfig
	TraceProvider  *trace.TracerProvider
	MeterProvider  *metric.MeterProvider
	LoggerProvider *log.LoggerProvider
	shutdown       func(context.Context) error
}

func NewTelemetryManager(
	lc fx.Lifecycle,
	metadata *config.Metadata,
	telemetryCfg *config.TelemetryConfig,
) (*TelemetryManager, error) {
	m := &TelemetryManager{
		metadata:     metadata,
		telemetryCfg: telemetryCfg,
	}

	ctx := context.Background()
	res, err := m.createResource(ctx)
	if err != nil {
		return nil, err
	}

	shutdown, providers, err := telemetry.SetupOpenTelemetrySdk(ctx, res, m.telemetryCfg)
	if err != nil {
		return nil, fmt.Errorf("setup otel sdk: %w", err)
	}
	m.TraceProvider = providers.TraceProvider
	m.MeterProvider = providers.MeterProvider
	m.LoggerProvider = providers.LoggerProvider
	m.shutdown = shutdown

	if err := runtime.Start(runtime.WithMeterProvider(m.MeterProvider)); err != nil {
		return nil, fmt.Errorf("start runtime metrics: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return m.shutdown(ctx)
		},
	})

	return m, nil
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
