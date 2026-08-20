package support

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/opentelemetry"
	"go.temporal.io/sdk/worker"
	"go.uber.org/fx"

	"github.com/codejsha/bookstore-microservices/inventory/internal/config"
	wf "github.com/codejsha/bookstore-microservices/inventory/internal/domain/workflow"
)

const temporalMeterName = "bookstore.inventory.temporal"

type TemporalWorker struct {
	client client.Client
	worker worker.Worker
}

func NewTemporalWorker(
	lc fx.Lifecycle,
	cfg *config.TemporalConfig,
	activities *wf.StockActivities,
	telemetryManager *TelemetryManager,
) (*TemporalWorker, error) {
	metricsHandler := opentelemetry.NewMetricsHandler(opentelemetry.MetricsHandlerOptions{
		Meter:                telemetryManager.MeterProvider.Meter(temporalMeterName),
		UseMonotonicCounters: true,
	})

	c, err := client.Dial(client.Options{
		HostPort:       cfg.Host,
		Namespace:      cfg.Namespace,
		MetricsHandler: metricsHandler,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create temporal client: %w", err)
	}

	w := worker.New(c, cfg.TaskQueue, worker.Options{})
	w.RegisterActivity(activities.ReserveStock)
	w.RegisterActivity(activities.ReleaseStock)

	tw := &TemporalWorker{
		client: c,
		worker: w,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := tw.worker.Start(); err != nil {
				return fmt.Errorf("start temporal worker: %w", err)
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			tw.worker.Stop()
			tw.client.Close()
			return nil
		},
	})

	return tw, nil
}
