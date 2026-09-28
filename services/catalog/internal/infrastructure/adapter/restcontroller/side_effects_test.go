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
	if !WaitForSideEffects(10 * time.Second) {
		t.Fatal("earlier side effects did not drain")
	}
	pub := &concurrencyPublisher{}
	SetEventPublisher(pub)
	defer SetEventPublisher(nil)

	total := maxInFlightSideEffects * 3
	for i := 0; i < total; i++ {
		dispatchSideEffects(context.Background(), "resource", fmt.Sprintf("r-%d", i), "created", logrus.Fields{})
	}
	if !WaitForSideEffects(10 * time.Second) {
		t.Fatal("side effects did not drain")
	}

	if peak := pub.peak.Load(); peak > maxInFlightSideEffects {
		t.Errorf("peak concurrent side effects = %d, want <= %d", peak, maxInFlightSideEffects)
	}
	if calls := pub.calls.Load(); int(calls) != total {
		t.Errorf("published side effects = %d, want %d", calls, total)
	}
}
