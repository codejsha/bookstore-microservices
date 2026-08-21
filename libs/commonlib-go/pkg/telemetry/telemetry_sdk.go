package telemetry

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/resource"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/config"
)

func SetupOpenTelemetrySdk(
	ctx context.Context,
	res *resource.Resource,
	telemetryCfg *config.TelemetryConfig,
) (
	shutdown func(context.Context) error,
	providers *Providers,
	err error,
) {
	// shutdown functions
	var shutdownFuncs []func(context.Context) error
	shutdown = func(ctx context.Context) error {
		var combinedErr error
		for _, fn := range shutdownFuncs {
			combinedErr = errors.Join(combinedErr, fn(ctx))
		}
		shutdownFuncs = nil
		return combinedErr
	}
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// collector client connection
	conn, err := CreateCollectorClient(telemetryCfg)
	if err != nil {
		handleErr(err)
		return
	}

	// propagator
	prop := CreatePropagator()
	otel.SetTextMapPropagator(prop)

	// trace provider
	tracerProvider, err := CreateTraceProvider(ctx, res, conn)
	if err != nil {
		handleErr(err)
		return
	}
	otel.SetTracerProvider(tracerProvider)
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)

	// meter provider
	meterProvider, err := CreateMeterProvider(ctx, res, conn)
	if err != nil {
		handleErr(err)
		return
	}
	otel.SetMeterProvider(meterProvider)
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)

	// logger provider
	loggerProvider, err := CreateLoggerProvider(ctx, res, conn)
	if err != nil {
		handleErr(err)
		return
	}
	global.SetLoggerProvider(loggerProvider)
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)

	providers = &Providers{
		TraceProvider:  tracerProvider,
		MeterProvider:  meterProvider,
		LoggerProvider: loggerProvider,
	}
	return shutdown, providers, nil
}
