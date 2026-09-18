package support

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/fx"

	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/adapter/restcontroller"
)

const (
	shutdownDrainDelay       = 5 * time.Second
	drainStepBudget          = shutdownDrainDelay + time.Second
	httpShutdownTimeout      = writeTimeout
	sideEffectsDrainTimeout  = 10 * time.Second
	kafkaCloseTimeout        = 5 * time.Second
	grpcClientCloseTimeout   = 5 * time.Second
	readinessCloseTimeout    = 2 * time.Second
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
	publisher io.Closer,
	readinessDataSource *ReadinessDataSource,
	telemetryManager *TelemetryManager,
	grpcClients ...io.Closer,
) {
	steps := shutdownSequence(ginServer, publisher, readinessDataSource, telemetryManager, grpcClients...)
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			runShutdownSteps(steps)
			return nil
		},
	})
}

func shutdownSequence(
	ginServer *GinServer,
	publisher io.Closer,
	readinessDataSource *ReadinessDataSource,
	telemetryManager *TelemetryManager,
	grpcClients ...io.Closer,
) []shutdownStep {
	return []shutdownStep{
		{
			name:   "http-drain-delay",
			budget: drainStepBudget,
			run: func(ctx context.Context) error {
				ginServer.BeginDrain()
				timer := time.NewTimer(shutdownDrainDelay)
				defer timer.Stop()
				select {
				case <-timer.C:
					return nil
				case <-ctx.Done():
					return ctx.Err()
				}
			},
		},
		{
			name:   "http-shutdown",
			budget: httpShutdownTimeout,
			run:    ginServer.Shutdown,
		},
		{
			name:   "side-effects-drain",
			budget: sideEffectsDrainTimeout,
			run: func(ctx context.Context) error {
				restcontroller.WaitForSideEffects(ctx)
				return nil
			},
		},
		{
			name:   "kafka-publisher-close",
			budget: kafkaCloseTimeout,
			run: func(context.Context) error {
				return publisher.Close()
			},
		},
		{
			name:   "grpc-clients-close",
			budget: grpcClientCloseTimeout,
			run: func(context.Context) error {
				return closeAll(grpcClients)
			},
		},
		{
			name:   "readiness-db-close",
			budget: readinessCloseTimeout,
			run: func(context.Context) error {
				return readinessDataSource.Close()
			},
		},
		{
			name:   "telemetry-shutdown",
			budget: telemetryShutdownTimeout,
			run:    telemetryManager.Shutdown,
		},
	}
}

func runShutdownSteps(steps []shutdownStep) {
	for _, step := range steps {
		runShutdownStep(step)
	}
}

func runShutdownStep(step shutdownStep) {
	ctx, cancel := context.WithTimeout(context.Background(), step.budget)
	defer cancel()

	start := time.Now()
	done := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- fmt.Errorf("shutdown step %s panicked: %v", step.name, r)
			}
		}()
		done <- step.run(ctx)
	}()

	var err error
	select {
	case err = <-done:
	case <-ctx.Done():
		err = ctx.Err()
	}

	entry := logrus.WithField("step", step.name).
		WithField("budget_ms", step.budget.Milliseconds()).
		WithField("duration_ms", time.Since(start).Milliseconds())
	switch {
	case err == nil:
		entry.Info("shutdown step completed")
	case errors.Is(err, context.DeadlineExceeded):
		entry.WithError(err).Warn("shutdown step timed out")
	default:
		entry.WithError(err).Error("shutdown step failed")
	}
}

func closeAll(closers []io.Closer) error {
	errs := make([]error, len(closers))
	var wg sync.WaitGroup
	for i, closer := range closers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs[i] = closer.Close()
		}()
	}
	wg.Wait()
	return errors.Join(errs...)
}
