package support

import (
	"context"
	"fmt"
	"sync"
	"time"

	restclient "github.com/codejsha/shared-library-go/pkg/rest/client"
	"github.com/sirupsen/logrus"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/opentelemetry"
	"go.temporal.io/sdk/worker"
	"go.uber.org/fx"

	"github.com/codejsha/bookstore-microservices/inventory/internal/config"
	wf "github.com/codejsha/bookstore-microservices/inventory/internal/domain/workflow"
)

const (
	temporalMeterName = "bookstore.inventory.temporal"

	workerStopTimeout         = 10 * time.Second
	workerConnectRetryInitial = 1 * time.Second
	workerConnectRetryMax     = 30 * time.Second
)

func temporalWorkerOptions() worker.Options {
	return worker.Options{
		WorkerStopTimeout: workerStopTimeout,
	}
}

type TemporalWorker struct {
	client    client.Client
	newWorker func() worker.Worker

	mu       sync.Mutex
	worker   worker.Worker
	stopped  bool
	launched bool

	stopC     chan struct{}
	doneC     chan struct{}
	closeOnce sync.Once

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

	c, err := client.NewLazyClient(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporal client: %w", err)
	}

	tw := &TemporalWorker{
		client: c,
		newWorker: func() worker.Worker {
			w := worker.New(c, cfg.TaskQueue, temporalWorkerOptions())
			w.RegisterActivity(activities.ReserveStock)
			w.RegisterActivity(activities.ReleaseStock)
			return w
		},
		stopC:  make(chan struct{}),
		doneC:  make(chan struct{}),
		tokens: tokens,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if tw.tokens != nil {
				tw.tokens.Start()
			}
			tw.mu.Lock()
			tw.launched = true
			tw.mu.Unlock()
			go tw.connect()
			return nil
		},
	})

	return tw, nil
}

func (tw *TemporalWorker) StopWorker() {
	_ = tw.stopWorker(context.Background())
}

func (tw *TemporalWorker) CloseClient() {
	tw.closeOnce.Do(func() {
		tw.client.Close()
		if tw.tokens != nil {
			tw.tokens.Stop()
		}
	})
}

func (tw *TemporalWorker) connect() {
	defer close(tw.doneC)

	delay := workerConnectRetryInitial
	for {
		if tw.isStopped() {
			return
		}

		w := tw.newWorker()
		if err := w.Start(); err != nil {
			logrus.Warnf("temporal worker start failed, retrying in %s: %v", delay, err)
			select {
			case <-tw.stopC:
				return
			case <-time.After(delay):
			}
			delay *= 2
			if delay > workerConnectRetryMax {
				delay = workerConnectRetryMax
			}
			continue
		}

		if !tw.adopt(w) {
			w.Stop()
			return
		}
		logrus.Info("temporal worker started")
		return
	}
}

func (tw *TemporalWorker) adopt(w worker.Worker) bool {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.stopped {
		return false
	}
	tw.worker = w
	return true
}

func (tw *TemporalWorker) isStopped() bool {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	return tw.stopped
}

func (tw *TemporalWorker) stopWorker(ctx context.Context) error {
	tw.mu.Lock()
	if tw.stopped {
		tw.mu.Unlock()
		return nil
	}
	tw.stopped = true
	launched := tw.launched
	close(tw.stopC)
	tw.mu.Unlock()

	if launched {
		select {
		case <-tw.doneC:
		case <-ctx.Done():
		}
	}

	tw.mu.Lock()
	w := tw.worker
	tw.worker = nil
	tw.mu.Unlock()

	if w != nil {
		w.Stop()
	}
	return nil
}

func (tw *TemporalWorker) shutdown(ctx context.Context) error {
	if err := tw.stopWorker(ctx); err != nil {
		return err
	}
	tw.CloseClient()
	return nil
}
