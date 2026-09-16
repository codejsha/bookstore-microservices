package support

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/fx"

	"github.com/codejsha/shared-library-go/pkg/database"
	"github.com/codejsha/shared-library-go/pkg/message"

	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/adapter/restcontroller"
)

const (
	shutdownDrainDelay       = 5 * time.Second
	shutdownDrainStepTimeout = shutdownDrainDelay + time.Second
	httpShutdownTimeout      = writeTimeout
	sideEffectsDrainTimeout  = 10 * time.Second
	databaseCloseTimeout     = 3 * time.Second
	kafkaCloseTimeout        = 5 * time.Second
	cacheCloseTimeout        = 2 * time.Second
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
	pub *message.KafkaAsyncPublisher,
	cacheCloser *CacheClientCloser,
	dataSource *database.DataSource,
	readinessDataSource *ReadinessDataSource,
	telemetryManager *TelemetryManager,
) {
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			runShutdownSteps(shutdownSequence(ginServer, pub, cacheCloser, dataSource, readinessDataSource, telemetryManager))
			return nil
		},
	})
}

func shutdownSequence(
	ginServer *GinServer,
	pub *message.KafkaAsyncPublisher,
	cacheCloser *CacheClientCloser,
	dataSource *database.DataSource,
	readinessDataSource *ReadinessDataSource,
	telemetryManager *TelemetryManager,
) []shutdownStep {
	return []shutdownStep{
		{
			name:   "http-drain-delay",
			budget: shutdownDrainStepTimeout,
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
			run: func(context.Context) error {
				if !restcontroller.WaitForSideEffects(sideEffectsDrainTimeout) {
					return errors.New("side-effect goroutines still running")
				}
				return nil
			},
		},
		{
			name:   "database-close",
			budget: databaseCloseTimeout,
			run: func(context.Context) error {
				return closeDataSourcePool(dataSource)
			},
		},
		{
			name:   "kafka-close",
			budget: kafkaCloseTimeout,
			run: func(context.Context) error {
				return pub.Close()
			},
		},
		{
			name:   "valkey-close",
			budget: cacheCloseTimeout,
			run: func(context.Context) error {
				return cacheCloser.Close()
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
				done <- fmt.Errorf("panic: %v", r)
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

	entry := logrus.WithFields(logrus.Fields{
		"step":     step.name,
		"budget":   step.budget.String(),
		"duration": time.Since(start).String(),
	})
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		entry.WithError(err).Warn("shutdown step timed out")
	case err != nil:
		entry.WithError(err).Warn("shutdown step failed")
	default:
		entry.Info("shutdown step completed")
	}
}

func closeDataSourcePool(dataSource *database.DataSource) error {
	if dataSource == nil {
		return nil
	}
	sqlDB, err := dataSource.DB().DB()
	if err != nil {
		return fmt.Errorf("resolve database handle: %w", err)
	}
	return closeConnectionPool(sqlDB)
}

func closeConnectionPool(sqlDB *sql.DB) error {
	if sqlDB == nil {
		return nil
	}
	return sqlDB.Close()
}
