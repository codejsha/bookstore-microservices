package support

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

type recordedStep struct {
	name        string
	hasDeadline bool
	remaining   time.Duration
}

type stepRecorder struct {
	mu    sync.Mutex
	steps []recordedStep
}

func (r *stepRecorder) record(name string, ctx context.Context) {
	deadline, ok := ctx.Deadline()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.steps = append(r.steps, recordedStep{name: name, hasDeadline: ok, remaining: time.Until(deadline)})
}

func TestRunShutdownSteps_earlierStepFailsOrOverruns_runsAllStepsInOrder(t *testing.T) {
	const budget = 40 * time.Millisecond
	rec := &stepRecorder{}
	steps := []shutdownStep{
		{name: "fails", budget: budget, run: func(ctx context.Context) error {
			rec.record("fails", ctx)
			return errors.New("boom")
		}},
		{name: "overruns", budget: budget, run: func(ctx context.Context) error {
			rec.record("overruns", ctx)
			time.Sleep(3 * budget)
			return ctx.Err()
		}},
		{name: "waitsOnBlockedCall", budget: budget, run: func(ctx context.Context) error {
			rec.record("waitsOnBlockedCall", ctx)
			return waitWithContext(ctx, func() error {
				select {}
			})
		}},
		{name: "last", budget: budget, run: func(ctx context.Context) error {
			rec.record("last", ctx)
			return nil
		}},
	}

	runShutdownSteps(steps)

	want := []string{"fails", "overruns", "waitsOnBlockedCall", "last"}
	if len(rec.steps) != len(want) {
		t.Fatalf("ran %d steps, want %d", len(rec.steps), len(want))
	}
	for i, step := range rec.steps {
		if step.name != want[i] {
			t.Fatalf("step %d = %q, want %q", i, step.name, want[i])
		}
		if !step.hasDeadline {
			t.Fatalf("step %q has no deadline", step.name)
		}
		if step.remaining <= budget/2 || step.remaining > budget {
			t.Fatalf("step %q remaining deadline = %v, want within (%v, %v]", step.name, step.remaining, budget/2, budget)
		}
	}
}

func TestWaitWithContext_callBlocksPastDeadline_stopsWaiting(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	release := make(chan struct{})
	defer close(release)

	start := time.Now()
	err := waitWithContext(ctx, func() error {
		<-release
		return nil
	})

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want %v", err, context.DeadlineExceeded)
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("waited %v after deadline", elapsed)
	}
}

func TestConcurrentShutdownStep_parallelParts_waitsForAllWithOwnBudgets(t *testing.T) {
	const partDelay = 50 * time.Millisecond
	rec := &stepRecorder{}
	part := func(name string, budget time.Duration) shutdownStep {
		return shutdownStep{name: name, budget: budget, run: func(ctx context.Context) error {
			time.Sleep(partDelay)
			rec.record(name, ctx)
			return nil
		}}
	}
	step := concurrentShutdownStep("intake", part("short", 80*time.Millisecond), part("long", 200*time.Millisecond))
	var afterRan int

	start := time.Now()
	runShutdownSteps([]shutdownStep{
		step,
		{name: "after", budget: 20 * time.Millisecond, run: func(context.Context) error {
			rec.mu.Lock()
			afterRan = len(rec.steps)
			rec.mu.Unlock()
			return nil
		}},
	})
	elapsed := time.Since(start)

	if step.budget != 200*time.Millisecond {
		t.Fatalf("step budget = %v, want %v", step.budget, 200*time.Millisecond)
	}
	if afterRan != 2 {
		t.Fatalf("next step saw %d finished parts, want 2", afterRan)
	}
	if elapsed >= 2*partDelay {
		t.Fatalf("parts took %v, want concurrent execution under %v", elapsed, 2*partDelay)
	}
	for _, s := range rec.steps {
		limit := 80 * time.Millisecond
		if s.name == "long" {
			limit = 200 * time.Millisecond
		}
		if !s.hasDeadline || s.remaining > limit-partDelay {
			t.Fatalf("part %q remaining deadline = %v, want at most %v", s.name, s.remaining, limit-partDelay)
		}
	}
}
