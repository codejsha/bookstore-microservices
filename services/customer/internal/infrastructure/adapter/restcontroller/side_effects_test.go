package restcontroller

import (
	"context"
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

func TestDispatchSideEffects_RecoversPanic(t *testing.T) {
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

func TestWaitForSideEffects_DrainsInFlight(t *testing.T) {
	var ran atomic.Bool
	SetEventPublisher(&slowPublisher{ran: &ran})
	defer SetEventPublisher(nil)

	dispatchSideEffects(context.Background(), "point", "p-1", "earned", logrus.Fields{})
	WaitForSideEffects(context.Background())

	if !ran.Load() {
		t.Error("WaitForSideEffects returned before the in-flight side effect completed")
	}
}
