package restcontroller

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/codejsha/shared-library-go/pkg/message"
)

const sideEffectsTimeout = 5 * time.Second

var sideEffectsWG sync.WaitGroup
var eventPublisher atomic.Pointer[message.EventPublisher]

type CacheInvalidator interface {
	Invalidate(ctx context.Context) error
}

var cacheInvalidator atomic.Pointer[CacheInvalidator]

func SetEventPublisher(p message.EventPublisher) {
	if p == nil {
		eventPublisher.Store(nil)
		return
	}
	eventPublisher.Store(&p)
}

func SetCacheInvalidator(inv CacheInvalidator) {
	if inv == nil {
		cacheInvalidator.Store(nil)
		return
	}
	cacheInvalidator.Store(&inv)
}

type resourceEvent struct {
	IdValue        string                 `json:"id"`
	TypeValue      string                 `json:"type"`
	MetadataValue  map[string]interface{} `json:"metadata"`
	CreatedAtValue time.Time              `json:"created_at"`
}

func (e resourceEvent) Id() string                       { return e.IdValue }
func (e resourceEvent) Type() string                     { return e.TypeValue }
func (e resourceEvent) Metadata() map[string]interface{} { return e.MetadataValue }
func (e resourceEvent) CreatedAt() time.Time             { return e.CreatedAtValue }

func dispatchSideEffects(parent context.Context, resource, key, action string, extra logrus.Fields) {
	sideEffectsWG.Add(1)
	go func() {
		defer sideEffectsWG.Done()
		defer func() {
			if r := recover(); r != nil {
				logrus.WithField("resource", resource).
					WithField("key", key).
					WithField("action", action).
					WithField("panic", r).
					Error("side-effect goroutine panicked; recovered")
			}
		}()
		runSideEffects(parent, resource, key, action, extra)
	}()
}

func WaitForSideEffects(timeout time.Duration) bool {
	done := make(chan struct{})
	go func() {
		sideEffectsWG.Wait()
		close(done)
	}()
	select {
	case <-done:
		return true
	case <-time.After(timeout):
		return false
	}
}

func runSideEffects(parent context.Context, resource, key, action string, extra logrus.Fields) {
	ctx, cancel := context.WithTimeout(parent, sideEffectsTimeout)
	defer cancel()

	g, gctx := errgroup.WithContext(ctx)
	g.Go(guardSideEffect("event", resource, key, func() error { return emitEvent(gctx, resource, key, action, extra) }))
	g.Go(guardSideEffect("cache", resource, key, func() error { return invalidateCache(gctx, resource, key) }))
	g.Go(guardSideEffect("audit", resource, key, func() error { return recordAudit(gctx, resource, key, action, extra) }))

	if err := g.Wait(); err != nil {
		logrus.WithError(err).
			WithField("resource", resource).
			WithField("key", key).
			WithField("action", action).
			Warn("side-effect work partially failed")
	}
}

func guardSideEffect(name, resource, key string, fn func() error) func() error {
	return func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				logrus.WithField("side_effect", name).
					WithField("resource", resource).
					WithField("key", key).
					WithField("panic", r).
					Error("side-effect task panicked; recovered")
				err = fmt.Errorf("side-effect %s panicked: %v", name, r)
			}
		}()
		return fn()
	}
}

func emitEvent(ctx context.Context, resource, key, action string, extra logrus.Fields) error {
	logrus.WithContext(ctx).
		WithField("resource", resource).
		WithField("key", key).
		WithField("action", action).
		WithFields(extra).
		Info("event emitted")

	pub := eventPublisher.Load()
	if pub == nil {
		return nil
	}
	metadata := make(map[string]interface{}, len(extra)+1)
	for k, v := range extra {
		metadata[k] = v
	}
	metadata["action"] = action
	evt := resourceEvent{
		IdValue:        key,
		TypeValue:      resource + "." + action,
		MetadataValue:  metadata,
		CreatedAtValue: time.Now().UTC(),
	}
	return (*pub).Publish(ctx, resource+".changed.v1", evt)
}

func invalidateCache(ctx context.Context, resource, key string) error {
	inv := cacheInvalidator.Load()
	if inv == nil {
		logrus.WithContext(ctx).
			WithField("resource", resource).
			WithField("key", key).
			Debug("cache invalidation skipped (no cache configured)")
		return nil
	}
	if err := (*inv).Invalidate(ctx); err != nil {
		logrus.WithContext(ctx).WithError(err).
			WithField("resource", resource).
			WithField("key", key).
			Warn("cache epoch bump failed; relying on list TTL")
		return nil
	}
	logrus.WithContext(ctx).
		WithField("resource", resource).
		WithField("key", key).
		Debug("cache invalidated (epoch bumped)")
	return nil
}

func recordAudit(ctx context.Context, resource, key, action string, extra logrus.Fields) error {
	logrus.WithContext(ctx).
		WithField("resource", resource).
		WithField("key", key).
		WithField("action", action).
		WithFields(extra).
		Info("audit recorded")
	return nil
}
