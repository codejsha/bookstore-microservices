package restcontroller

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/codejsha/shared-library-go/pkg/message"
)

type panicPublisher struct{ called atomic.Bool }

func (p *panicPublisher) Publish(_ context.Context, _ string, _ message.EventMessage) error {
	p.called.Store(true)
	panic("boom in publish")
}

func TestDispatchSideEffects_WhenSideEffectPanics_RecoversAndReleasesWaiter(t *testing.T) {
	pub := &panicPublisher{}
	SetEventPublisher(pub)
	defer SetEventPublisher(nil)

	dispatchSideEffects(context.Background(), "review", "r-1", "created", logrus.Fields{})

	done := make(chan struct{})
	go func() {
		WaitForSideEffects(context.Background())
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("WaitForSideEffects did not return after a panicking side effect")
	}
	if !pub.called.Load() {
		t.Error("publisher was never invoked")
	}
}

type slowPublisher struct{ ran *atomic.Bool }

func (s *slowPublisher) Publish(_ context.Context, _ string, _ message.EventMessage) error {
	time.Sleep(50 * time.Millisecond)
	s.ran.Store(true)
	return nil
}

func TestWaitForSideEffects_WhenSideEffectInFlight_BlocksUntilItCompletes(t *testing.T) {
	var ran atomic.Bool
	SetEventPublisher(&slowPublisher{ran: &ran})
	defer SetEventPublisher(nil)

	dispatchSideEffects(context.Background(), "point", "p-1", "earned", logrus.Fields{})
	WaitForSideEffects(context.Background())

	if !ran.Load() {
		t.Error("WaitForSideEffects returned before the in-flight side effect completed")
	}
}

type concurrencyPublisher struct {
	inFlight atomic.Int32
	peak     atomic.Int32
	calls    atomic.Int32
}

func (c *concurrencyPublisher) Publish(_ context.Context, _ string, _ message.EventMessage) error {
	n := c.inFlight.Add(1)
	defer c.inFlight.Add(-1)
	for {
		peak := c.peak.Load()
		if n <= peak || c.peak.CompareAndSwap(peak, n) {
			break
		}
	}
	time.Sleep(20 * time.Millisecond)
	c.calls.Add(1)
	return nil
}

func TestDispatchSideEffects_burstExceedsInFlightLimit_capsConcurrencyAndRunsAll(t *testing.T) {
	WaitForSideEffects(context.Background())
	pub := &concurrencyPublisher{}
	SetEventPublisher(pub)
	defer SetEventPublisher(nil)

	total := maxInFlightSideEffects * 3
	for i := 0; i < total; i++ {
		dispatchSideEffects(context.Background(), "resource", fmt.Sprintf("r-%d", i), "created", logrus.Fields{})
	}
	WaitForSideEffects(context.Background())

	if peak := pub.peak.Load(); peak > maxInFlightSideEffects {
		t.Errorf("peak concurrent side effects = %d, want <= %d", peak, maxInFlightSideEffects)
	}
	if calls := pub.calls.Load(); int(calls) != total {
		t.Errorf("published side effects = %d, want %d", calls, total)
	}
}
