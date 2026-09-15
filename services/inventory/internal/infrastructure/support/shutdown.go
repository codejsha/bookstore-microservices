package support

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/fx"

	"github.com/codejsha/shared-library-go/pkg/message"
)

const (
	shutdownDrainDelay       = 5 * time.Second
	httpShutdownTimeout      = writeTimeout
	temporalStopTimeout      = 10 * time.Second
	temporalCloseTimeout     = 5 * time.Second
	kafkaCloseTimeout        = 5 * time.Second
	telemetryShutdownTimeout = 5 * time.Second
	ShutdownStopTimeout      = 55 * time.Second
)

type shutdownStep struct {
	name   string
	budget time.Duration
	run    func(ctx context.Context) error
}

func RegisterShutdownSequence(
	lc fx.Lifecycle,
	ginServer *GinServer,
	temporalWorker *TemporalWorker,
	kafkaPublisher *message.KafkaAsyncPublisher,
	telemetryManager *TelemetryManager,
) {
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			runShutdownSteps([]shutdownStep{
				{
					name:   "drain",
					budget: shutdownDrainDelay,
					run: func(ctx context.Context) error {
						ginServer.BeginDrain()
						return sleepWithContext(ctx, shutdownDrainDelay)
					},
				},
				concurrentShutdownStep("stop-intake",
					shutdownStep{
						name:   "http-server",
						budget: httpShutdownTimeout,
						run:    ginServer.Shutdown,
					},
					shutdownStep{
						name:   "temporal-worker",
						budget: temporalStopTimeout,
						run: func(ctx context.Context) error {
							return waitWithContext(ctx, func() error {
								temporalWorker.StopWorker()
								return nil
							})
						},
					},
				),
				{
					name:   "kafka-publisher",
					budget: kafkaCloseTimeout,
					run: func(ctx context.Context) error {
						return waitWithContext(ctx, kafkaPublisher.Close)
					},
				},
				{
					name:   "temporal-client",
					budget: temporalCloseTimeout,
					run: func(ctx context.Context) error {
						return waitWithContext(ctx, func() error {
							temporalWorker.CloseClient()
							return nil
						})
					},
				},
				{
					name:   "telemetry",
					budget: telemetryShutdownTimeout,
					run:    telemetryManager.Shutdown,
				},
			})
			return nil
		},
	})
}

func runShutdownSteps(steps []shutdownStep) {
	for _, step := range steps {
		_ = runShutdownStep(context.Background(), step)
	}
}

func runShutdownStep(parent context.Context, step shutdownStep) error {
	ctx, cancel := context.WithTimeout(parent, step.budget)
	defer cancel()

	start := time.Now()
	err := step.run(ctx)
	entry := logrus.WithFields(logrus.Fields{
		"shutdown_step": step.name,
		"budget_ms":     step.budget.Milliseconds(),
		"duration_ms":   time.Since(start).Milliseconds(),
	})
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		entry.WithError(err).Warn("shutdown step timed out")
	case err != nil:
		entry.WithError(err).Error("shutdown step failed")
	default:
		entry.Info("shutdown step completed")
	}
	return err
}

func concurrentShutdownStep(name string, parts ...shutdownStep) shutdownStep {
	var budget time.Duration
	for _, part := range parts {
		budget = max(budget, part.budget)
	}
	return shutdownStep{
		name:   name,
		budget: budget,
		run: func(ctx context.Context) error {
			errs := make([]error, len(parts))
			var wg sync.WaitGroup
			for i, part := range parts {
				wg.Go(func() {
					errs[i] = runShutdownStep(ctx, part)
				})
			}
			wg.Wait()
			return errors.Join(errs...)
		},
	}
}

func waitWithContext(ctx context.Context, fn func() error) error {
	done := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- fmt.Errorf("panic: %v", r)
			}
		}()
		done <- fn()
	}()
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func sleepWithContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-timer.C:
	case <-ctx.Done():
	}
	return nil
}
