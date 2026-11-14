package restcontroller

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/codejsha/shared-library-go/pkg/message"
)

type publishedEvent struct {
	topic    string
	id       string
	kind     string
	metadata map[string]interface{}
	at       time.Time
}

type recordingPublisher struct {
	mu     sync.Mutex
	events []publishedEvent
	err    error
	panics bool
}

func (p *recordingPublisher) Publish(_ context.Context, topic string, msg message.EventMessage) error {
	p.mu.Lock()
	p.events = append(p.events, publishedEvent{
		topic: topic, id: msg.Id(), kind: msg.Type(), metadata: msg.Metadata(), at: msg.CreatedAt(),
	})
	panics, err := p.panics, p.err
	p.mu.Unlock()
	if panics {
		panic("broker client exploded")
	}
	return err
}

func (p *recordingPublisher) published() []publishedEvent {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]publishedEvent(nil), p.events...)
}

type countingInvalidator struct {
	calls  atomic.Int32
	err    error
	panics bool
}

func (c *countingInvalidator) Invalidate(context.Context) error {
	c.calls.Add(1)
	if c.panics {
		panic("cache client exploded")
	}
	return c.err
}

func drainSideEffects(t *testing.T) {
	t.Helper()
	if !WaitForSideEffects(10 * time.Second) {
		t.Fatal("side effects did not drain")
	}
}

func captureSideEffects(t *testing.T) (*recordingPublisher, *countingInvalidator) {
	t.Helper()
	drainSideEffects(t)
	pub := &recordingPublisher{}
	inv := &countingInvalidator{}
	SetEventPublisher(pub)
	SetCacheInvalidator(inv)
	t.Cleanup(func() {
		WaitForSideEffects(10 * time.Second)
		SetEventPublisher(nil)
		SetCacheInvalidator(nil)
	})
	return pub, inv
}

func TestDispatchSideEffects_resourceChanged_publishesEventAndBumpsCacheOnce(t *testing.T) {
	pub, inv := captureSideEffects(t)
	before := time.Now().UTC()

	dispatchSideEffects(context.Background(), "work", "w-1", "created", logrus.Fields{"title": "Dune"})
	drainSideEffects(t)

	events := pub.published()
	if len(events) != 1 {
		t.Fatalf("published events = %d, want 1", len(events))
	}
	evt := events[0]
	if evt.topic != "work.changed.v1" || evt.id != "w-1" || evt.kind != "work.created" {
		t.Errorf("event = %+v, want topic work.changed.v1, id w-1, type work.created", evt)
	}
	if evt.metadata["title"] != "Dune" || evt.metadata["action"] != "created" || len(evt.metadata) != 2 {
		t.Errorf("metadata = %v, want title and action only", evt.metadata)
	}
	if evt.at.Before(before) || evt.at.Location() != time.UTC {
		t.Errorf("created_at = %v, want a UTC time not before %v", evt.at, before)
	}
	if got := inv.calls.Load(); got != 1 {
		t.Errorf("cache invalidations = %d, want 1", got)
	}
}

func TestDispatchSideEffects_cacheInvalidationFails_stillPublishesEvent(t *testing.T) {
	pub, inv := captureSideEffects(t)
	inv.err = errors.New("valkey unreachable")

	dispatchSideEffects(context.Background(), "author", "a-1", "updated", logrus.Fields{})
	drainSideEffects(t)

	if events := pub.published(); len(events) != 1 || events[0].kind != "author.updated" {
		t.Errorf("published events = %+v, want one author.updated", events)
	}
	if got := inv.calls.Load(); got != 1 {
		t.Errorf("cache invalidations = %d, want 1 attempt", got)
	}
}

func TestDispatchSideEffects_publisherFails_stillAttemptsCacheInvalidation(t *testing.T) {
	pub, inv := captureSideEffects(t)
	pub.err = errors.New("broker unreachable")

	dispatchSideEffects(context.Background(), "work", "w-1", "updated", logrus.Fields{})
	drainSideEffects(t)

	if got := len(pub.published()); got != 1 {
		t.Errorf("publish attempts = %d, want 1", got)
	}
	if got := inv.calls.Load(); got != 1 {
		t.Errorf("cache invalidations = %d, want 1", got)
	}
}

func TestDispatchSideEffects_publisherPanics_recoversAndKeepsServingLaterDispatches(t *testing.T) {
	pub, inv := captureSideEffects(t)
	pub.panics = true

	dispatchSideEffects(context.Background(), "work", "w-1", "created", logrus.Fields{})
	drainSideEffects(t)
	pub.mu.Lock()
	pub.panics = false
	pub.mu.Unlock()
	dispatchSideEffects(context.Background(), "work", "w-2", "created", logrus.Fields{})
	drainSideEffects(t)

	events := pub.published()
	if len(events) != 2 || events[1].id != "w-2" {
		t.Errorf("published events = %+v, want the second dispatch delivered after the panic", events)
	}
	if got := inv.calls.Load(); got != 2 {
		t.Errorf("cache invalidations = %d, want 2", got)
	}
}

func TestDispatchSideEffects_invalidatorPanics_recoversAndStillPublishes(t *testing.T) {
	pub, inv := captureSideEffects(t)
	inv.panics = true

	dispatchSideEffects(context.Background(), "subject", "s-1", "created", logrus.Fields{"name": "Fantasy"})
	drainSideEffects(t)

	if events := pub.published(); len(events) != 1 || events[0].topic != "subject.changed.v1" {
		t.Errorf("published events = %+v, want one on subject.changed.v1", events)
	}
}

func TestDispatchSideEffects_parentContextAlreadyCancelled_dropsWork(t *testing.T) {
	pub, inv := captureSideEffects(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	dispatchSideEffects(ctx, "work", "w-1", "created", logrus.Fields{})
	drainSideEffects(t)

	if events := pub.published(); len(events) != 0 {
		t.Errorf("published events = %+v, want none", events)
	}
	if got := inv.calls.Load(); got != 0 {
		t.Errorf("cache invalidations = %d, want 0", got)
	}
}

func TestDispatchSideEffects_noPublisherOrInvalidatorConfigured_completesWithoutThem(t *testing.T) {
	pub, inv := captureSideEffects(t)
	SetEventPublisher(nil)
	SetCacheInvalidator(nil)

	dispatchSideEffects(context.Background(), "work", "w-1", "created", logrus.Fields{})

	if !WaitForSideEffects(10 * time.Second) {
		t.Fatal("side effects did not drain")
	}
	if len(pub.published()) != 0 || inv.calls.Load() != 0 {
		t.Errorf("detached publisher/invalidator were still called: %d events, %d invalidations", len(pub.published()), inv.calls.Load())
	}
}
