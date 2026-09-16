package support

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/fx"

	"github.com/codejsha/shared-library-go/pkg/database"
	"github.com/codejsha/shared-library-go/pkg/message"

	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure/adapter/restcontroller"
)

const (
	shutdownDrainDelay       = 5 * time.Second
	httpShutdownTimeout      = writeTimeout
	grpcShutdownTimeout      = 10 * time.Second
	sideEffectsDrainTimeout  = 10 * time.Second
	databaseCloseTimeout     = 3 * time.Second
	kafkaCloseTimeout        = 5 * time.Second
	cacheCloseTimeout        = 5 * time.Second
	readinessCloseTimeout    = 2 * time.Second
	telemetryShutdownTimeout = 5 * time.Second
	ShutdownStopTimeout      = 55 * time.Second
)

type shutdownStep struct {
	name     string
	budget   time.Duration
	run      func(ctx context.Context) error
	parallel []shutdownStep
}

func RegisterShutdownSequence(
	lc fx.Lifecycle,
	ginServer *GinServer,
	grpcServer *GrpcServer,
	publisher *message.KafkaAsyncPublisher,
	cacheClient *CacheClient,
	dataSource *database.DataSource,
	readinessDataSource *ReadinessDataSource,
	telemetryManager *TelemetryManager,
) {
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			runShutdownSteps([]shutdownStep{
				{
					name:   "drain readiness",
					budget: shutdownDrainDelay + time.Second,
					run: func(context.Context) error {
						ginServer.BeginDrain()
						time.Sleep(shutdownDrainDelay)
						return nil
					},
				},
				{
					name: "stop intake",
					parallel: []shutdownStep{
						{name: "http server shutdown", budget: httpShutdownTimeout, run: ginServer.Shutdown},
						{name: "grpc server shutdown", budget: grpcShutdownTimeout, run: grpcServer.Shutdown},
					},
				},
				{
					name:   "side effects drain",
					budget: sideEffectsDrainTimeout,
					run: func(ctx context.Context) error {
						restcontroller.WaitForSideEffects(ctx)
						return nil
					},
				},
				{
					name:   "database close",
					budget: databaseCloseTimeout,
					run: func(context.Context) error {
						return closeDataSourcePool(dataSource)
					},
				},
				{
					name:   "kafka publisher close",
					budget: kafkaCloseTimeout,
					run:    func(context.Context) error { return publisher.Close() },
				},
				{
					name:   "valkey client close",
					budget: cacheCloseTimeout,
					run:    func(context.Context) error { return cacheClient.Close() },
				},
				{
					name:   "readiness db close",
					budget: readinessCloseTimeout,
					run:    func(context.Context) error { return readinessDataSource.Close() },
				},
				{
					name:   "telemetry shutdown",
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
		if len(step.parallel) > 0 {
			runParallelShutdownSteps(step)
			continue
		}
		runShutdownStep(step)
	}
}

func runParallelShutdownSteps(group shutdownStep) {
	start := time.Now()
	var wg sync.WaitGroup
	for _, step := range group.parallel {
		wg.Add(1)
		go func(step shutdownStep) {
			defer wg.Done()
			runShutdownStep(step)
		}(step)
	}
	wg.Wait()
	logrus.WithFields(logrus.Fields{
		"step":        group.name,
		"duration_ms": time.Since(start).Milliseconds(),
	}).Info("shutdown step completed")
}

func runShutdownStep(step shutdownStep) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), step.budget)
	defer cancel()

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
		select {
		case err = <-done:
		default:
			err = ctx.Err()
		}
	}

	entry := logrus.WithFields(logrus.Fields{
		"step":        step.name,
		"budget_ms":   step.budget.Milliseconds(),
		"duration_ms": time.Since(start).Milliseconds(),
	})
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		entry.WithError(err).Warn("shutdown step timed out")
	case err != nil:
		entry.WithError(err).Error("shutdown step failed")
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
