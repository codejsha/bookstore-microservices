package pgsql

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestOptimisticRetryDelay_anyAttempt_staysWithinCap(t *testing.T) {
	for attempt := 0; attempt < 10; attempt++ {
		for sample := 0; sample < 100; sample++ {
			delay := optimisticRetryDelay(attempt)
			if delay <= 0 || delay > optimisticRetryMaxDelay {
				t.Fatalf("attempt %d delay = %v, want within (0, %v]", attempt, delay, optimisticRetryMaxDelay)
			}
		}
	}
}

func TestOptimisticRetryDelay_laterAttempts_backOffFurther(t *testing.T) {
	var firstMax, lateMin time.Duration
	lateMin = optimisticRetryMaxDelay
	for sample := 0; sample < 100; sample++ {
		if d := optimisticRetryDelay(0); d > firstMax {
			firstMax = d
		}
		if d := optimisticRetryDelay(5); d < lateMin {
			lateMin = d
		}
	}
	if lateMin <= firstMax {
		t.Fatalf("late attempt min delay %v must exceed first attempt max delay %v", lateMin, firstMax)
	}
}

func TestOptimisticRetryDelay_contendingWriters_jittered(t *testing.T) {
	seen := make(map[time.Duration]struct{})
	for sample := 0; sample < 50; sample++ {
		seen[optimisticRetryDelay(3)] = struct{}{}
	}
	if len(seen) < 2 {
		t.Fatal("delays must be jittered so contending writers do not retry in lockstep")
	}
}

func TestWaitBeforeOptimisticRetry_cancelledContext_abortsWithoutSleeping(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	start := time.Now()
	err := waitBeforeOptimisticRetry(ctx, 9)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled", err)
	}
	if elapsed := time.Since(start); elapsed >= optimisticRetryMaxDelay {
		t.Fatalf("waited %v, want an immediate return on a cancelled request", elapsed)
	}
}

func TestWaitBeforeOptimisticRetry_liveContext_sleepsThenRetries(t *testing.T) {
	if err := waitBeforeOptimisticRetry(context.Background(), 0); err != nil {
		t.Fatalf("wait: %v", err)
	}
}
