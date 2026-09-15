package support

import (
	"context"
	"errors"
	"slices"
	"sync"
	"testing"
	"time"
)

type stepRecorder struct {
	mu    sync.Mutex
	names []string
}

func (r *stepRecorder) record(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.names = append(r.names, name)
}

func (r *stepRecorder) snapshot() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return slices.Clone(r.names)
}

func TestRunShutdownSteps_earlierStepsFailOrOverrun_runsEveryStepInDeclaredOrder(t *testing.T) {
	rec := &stepRecorder{}
	release := make(chan struct{})
	t.Cleanup(func() { close(release) })

	runShutdownSteps([]shutdownStep{
		{name: "fails", budget: 20 * time.Millisecond, run: func(context.Context) error {
			rec.record("fails")
			return errors.New("boom")
		}},
		{name: "ignores-deadline", budget: 20 * time.Millisecond, run: func(context.Context) error {
			rec.record("ignores-deadline")
			<-release
			return nil
		}},
		{name: "honours-deadline", budget: 20 * time.Millisecond, run: func(ctx context.Context) error {
			rec.record("honours-deadline")
			<-ctx.Done()
			return ctx.Err()
		}},
		{name: "panics", budget: 20 * time.Millisecond, run: func(context.Context) error {
			rec.record("panics")
			panic("boom")
		}},
		{name: "last", budget: 20 * time.Millisecond, run: func(context.Context) error {
			rec.record("last")
			return nil
		}},
	})

	want := []string{"fails", "ignores-deadline", "honours-deadline", "panics", "last"}
	if got := rec.snapshot(); !slices.Equal(got, want) {
		t.Fatalf("step order = %v, want %v", got, want)
	}
}

func TestRunShutdownSteps_slowEarlierStep_laterStepKeepsFullBudget(t *testing.T) {
	const firstBudget = 50 * time.Millisecond
	const secondBudget = 80 * time.Millisecond
	var remaining time.Duration
	var hasDeadline bool

	runShutdownSteps([]shutdownStep{
		{name: "slow", budget: firstBudget, run: func(ctx context.Context) error {
			<-ctx.Done()
			return ctx.Err()
		}},
		{name: "measured", budget: secondBudget, run: func(ctx context.Context) error {
			var deadline time.Time
			deadline, hasDeadline = ctx.Deadline()
			remaining = time.Until(deadline)
			return nil
		}},
	})

	if !hasDeadline {
		t.Fatal("step context has no deadline")
	}
	if remaining <= secondBudget-firstBudget || remaining > secondBudget {
		t.Fatalf("remaining budget = %v, want within (%v, %v]", remaining, secondBudget-firstBudget, secondBudget)
	}
}

func TestShutdownSequence_catalogComponents_stopsInDependencyOrder(t *testing.T) {
	steps := shutdownSequence(nil, nil, nil, nil)

	names := make([]string, 0, len(steps))
	var total time.Duration
	for _, step := range steps {
		names = append(names, step.name)
		total += step.budget
	}

	want := []string{"http-drain-delay", "http-shutdown", "side-effects-drain", "kafka-close", "valkey-close", "telemetry-shutdown"}
	if !slices.Equal(names, want) {
		t.Fatalf("shutdown order = %v, want %v", names, want)
	}
	if total >= ShutdownStopTimeout {
		t.Fatalf("summed step budgets %v must stay under stop timeout %v", total, ShutdownStopTimeout)
	}
}

func TestCacheClientCloser_cacheDisabled_closesAsNoOp(t *testing.T) {
	if err := (&CacheClientCloser{}).Close(); err != nil {
		t.Fatalf("Close() = %v, want nil", err)
	}
	if err := (*CacheClientCloser)(nil).Close(); err != nil {
		t.Fatalf("nil Close() = %v, want nil", err)
	}
}
