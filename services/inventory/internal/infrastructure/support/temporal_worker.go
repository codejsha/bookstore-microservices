package support

import (
	"context"
	"fmt"

	restclient "github.com/codejsha/shared-library-go/pkg/rest/client"
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
	tokens *TemporalTokenSource
}

func NewTemporalWorker(
	lc fx.Lifecycle,
	cfg *config.TemporalConfig,
	activities *wf.StockActivities,
	telemetryManager *TelemetryManager,
	restyClient *restclient.RestyClient,
) (*TemporalWorker, error) {
	metricsHandler := opentelemetry.NewMetricsHandler(opentelemetry.MetricsHandlerOptions{
		Meter:                telemetryManager.MeterProvider.Meter(temporalMeterName),
		UseMonotonicCounters: true,
	})

	options := client.Options{
		HostPort:       cfg.Host,
		Namespace:      cfg.Namespace,
		MetricsHandler: metricsHandler,
	}

	var tokens *TemporalTokenSource
	if cfg.Auth != nil && cfg.Auth.Enabled {
		var err error
		tokens, err = NewTemporalTokenSource(cfg.Auth, restyClient)
		if err != nil {
			return nil, fmt.Errorf("invalid temporal auth config: %w", err)
		}
		initCtx, cancel := context.WithTimeout(context.Background(), temporalTokenInitTimeout)
		err = tokens.Init(initCtx)
		cancel()
		if err != nil {
			return nil, fmt.Errorf("failed to initialize temporal access token: %w", err)
		}
		options.HeadersProvider = tokens
	}

	c, err := client.Dial(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporal client: %w", err)
	}

	w := worker.New(c, cfg.TaskQueue, worker.Options{})
	w.RegisterActivity(activities.ReserveStock)
	w.RegisterActivity(activities.ReleaseStock)

	tw := &TemporalWorker{
		client: c,
		worker: w,
		tokens: tokens,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if tw.tokens != nil {
				tw.tokens.Start()
			}
			if err := tw.worker.Start(); err != nil {
				return fmt.Errorf("start temporal worker: %w", err)
			}
			return nil
		},
	})

	return tw, nil
}

func (tw *TemporalWorker) StopWorker() {
	tw.worker.Stop()
}

func (tw *TemporalWorker) CloseClient() {
	tw.client.Close()
	if tw.tokens != nil {
		tw.tokens.Stop()
	}
}
