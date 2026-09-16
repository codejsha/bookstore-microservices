package support

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/fx/fxtest"

	"github.com/codejsha/bookstore-microservices/inventory/internal/config"
	wf "github.com/codejsha/bookstore-microservices/inventory/internal/domain/workflow"
)

const unreachableTemporalHost = "127.0.0.1:1"

func newTestTemporalWorker(t *testing.T, lc *fxtest.Lifecycle) *TemporalWorker {
	t.Helper()
	tw, err := NewTemporalWorker(
		lc,
		&config.TemporalConfig{
			Host:      unreachableTemporalHost,
			Namespace: "bookstore",
			TaskQueue: "inventory-task-queue",
		},
		wf.NewStockActivities(nil),
		&TelemetryManager{MeterProvider: metric.NewMeterProvider()},
		nil,
	)
	if err != nil {
		t.Fatalf("NewTemporalWorker: %v", err)
	}
	return tw
}

func TestTemporalWorkerOptions_gracefulStop_setsNonZeroStopTimeout(t *testing.T) {
	opts := temporalWorkerOptions()

	if opts.WorkerStopTimeout <= 0 {
		t.Fatalf("WorkerStopTimeout = %v, want > 0 so in-flight activities get a grace period", opts.WorkerStopTimeout)
	}
	if opts.WorkerStopTimeout != workerStopTimeout {
		t.Fatalf("WorkerStopTimeout = %v, want %v", opts.WorkerStopTimeout, workerStopTimeout)
	}
}

func TestNewTemporalWorker_temporalUnreachable_constructorSucceeds(t *testing.T) {
	lc := fxtest.NewLifecycle(t)

	tw := newTestTemporalWorker(t, lc)

	if tw.client == nil {
		t.Fatal("client must be created lazily instead of dialled eagerly")
	}
}

func TestNewTemporalWorker_temporalUnreachable_appStartupSucceeds(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	tw := newTestTemporalWorker(t, lc)

	lc.RequireStart()
	t.Cleanup(func() { _ = tw.shutdown(context.Background()) })

	tw.mu.Lock()
	adopted := tw.worker != nil
	tw.mu.Unlock()
	if adopted {
		t.Fatal("worker must not be adopted while temporal is unreachable")
	}
}

func TestNewTemporalWorker_temporalUnreachable_shutdownStopsRetryLoop(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	tw := newTestTemporalWorker(t, lc)
	lc.RequireStart()

	if err := tw.shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown: %v", err)
	}

	select {
	case <-tw.doneC:
	case <-time.After(workerStopTimeout):
		t.Fatal("background connect loop still running after shutdown")
	}
}

func TestShutdown_calledTwice_staysIdempotent(t *testing.T) {
	lc := fxtest.NewLifecycle(t)
	tw := newTestTemporalWorker(t, lc)
	lc.RequireStart()

	if err := tw.shutdown(context.Background()); err != nil {
		t.Fatalf("first shutdown: %v", err)
	}
	if err := tw.shutdown(context.Background()); err != nil {
		t.Fatalf("second shutdown: %v", err)
	}
}
