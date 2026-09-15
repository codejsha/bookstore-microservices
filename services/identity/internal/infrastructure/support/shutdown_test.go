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

func TestRunShutdownSteps_earlierStepsFailOrOverrun_runsEveryStepInOrder(t *testing.T) {
	rec := &stepRecorder{}
	release := make(chan struct{})
	defer close(release)

	start := time.Now()
	runShutdownSteps([]shutdownStep{
		{name: "fails", budget: 20 * time.Millisecond, run: func(context.Context) error {
			rec.record("fails")
			return errors.New("close failed")
		}},
		{name: "ignores deadline", budget: 20 * time.Millisecond, run: func(context.Context) error {
			rec.record("ignores deadline")
			<-release
			return nil
		}},
		{name: "honors deadline", budget: 20 * time.Millisecond, run: func(ctx context.Context) error {
			rec.record("honors deadline")
			<-ctx.Done()
			return ctx.Err()
		}},
		{name: "panics", budget: 20 * time.Millisecond, run: func(context.Context) error {
			rec.record("panics")
			panic("boom")
		}},
		{name: "succeeds", budget: 20 * time.Millisecond, run: func(context.Context) error {
			rec.record("succeeds")
			return nil
		}},
	})
	elapsed := time.Since(start)

	want := []string{"fails", "ignores deadline", "honors deadline", "panics", "succeeds"}
	if got := rec.snapshot(); !slices.Equal(got, want) {
		t.Fatalf("step order = %v, want %v", got, want)
	}
	if elapsed > time.Second {
		t.Fatalf("sequence took %v, want bounded by step budgets", elapsed)
	}
}

func TestRunShutdownSteps_sequentialSteps_eachReceivesOwnDeadline(t *testing.T) {
	const firstBudget = 40 * time.Millisecond
	const secondBudget = 80 * time.Millisecond
	type observed struct {
		remaining   time.Duration
		hasDeadline bool
	}
	first := make(chan observed, 1)
	second := make(chan observed, 1)

	runShutdownSteps([]shutdownStep{
		{name: "first", budget: firstBudget, run: func(ctx context.Context) error {
			deadline, ok := ctx.Deadline()
			first <- observed{remaining: time.Until(deadline), hasDeadline: ok}
			<-ctx.Done()
			return ctx.Err()
		}},
		{name: "second", budget: secondBudget, run: func(ctx context.Context) error {
			deadline, ok := ctx.Deadline()
			second <- observed{remaining: time.Until(deadline), hasDeadline: ok}
			return nil
		}},
	})

	f, sec := <-first, <-second
	if !f.hasDeadline || f.remaining <= 0 || f.remaining > firstBudget {
		t.Fatalf("first step deadline = %v (set %v), want within (0, %v]", f.remaining, f.hasDeadline, firstBudget)
	}
	if !sec.hasDeadline || sec.remaining <= firstBudget || sec.remaining > secondBudget {
		t.Fatalf("second step deadline = %v (set %v), want within (%v, %v]", sec.remaining, sec.hasDeadline, firstBudget, secondBudget)
	}
}

func TestRunShutdownSteps_parallelGroup_runsMembersConcurrentlyBeforeNextStep(t *testing.T) {
	rec := &stepRecorder{}
	aStarted := make(chan struct{})
	bStarted := make(chan struct{})
	member := func(name string, own, other chan struct{}) shutdownStep {
		return shutdownStep{name: name, budget: time.Second, run: func(ctx context.Context) error {
			close(own)
			select {
			case <-other:
				rec.record(name)
				return nil
			case <-ctx.Done():
				rec.record(name + " alone")
				return ctx.Err()
			}
		}}
	}

	runShutdownSteps([]shutdownStep{
		{name: "group", parallel: []shutdownStep{
			member("a", aStarted, bStarted),
			member("b", bStarted, aStarted),
		}},
		{name: "next", budget: 20 * time.Millisecond, run: func(context.Context) error {
			rec.record("next")
			return nil
		}},
	})

	got := rec.snapshot()
	if len(got) != 3 || got[2] != "next" {
		t.Fatalf("step order = %v, want both group members then next", got)
	}
	members := slices.Sorted(slices.Values(got[:2]))
	if !slices.Equal(members, []string{"a", "b"}) {
		t.Fatalf("group members = %v, want a and b overlapping", members)
	}
}
