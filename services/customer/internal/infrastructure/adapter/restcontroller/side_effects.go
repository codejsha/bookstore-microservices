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

const (
	sideEffectTimeout       = 5 * time.Second
	sideEffectsDrainTimeout = 10 * time.Second
)

var eventPublisher atomic.Pointer[message.EventPublisher]

var sideEffectsWG sync.WaitGroup

func dispatchSideEffects(parent context.Context, resource, key, action string, extra logrus.Fields) {
	sideEffectsWG.Add(1)
	go func() {
		defer sideEffectsWG.Done()
		defer func() {
			if r := recover(); r != nil {
				logrus.WithField("panic", r).
					WithField("resource", resource).
					WithField("key", key).
					WithField("action", action).
					Error("side effects panicked; recovered")
			}
		}()
		runSideEffects(parent, resource, key, action, extra)
	}()
}

func WaitForSideEffects(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, sideEffectsDrainTimeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		sideEffectsWG.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		logrus.Warn("side-effects drain timed out; some in-flight events may be dropped")
	}
}

func SetEventPublisher(p message.EventPublisher) {
	if p == nil {
		eventPublisher.Store(nil)
		return
	}
	eventPublisher.Store(&p)
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

func runSideEffects(parent context.Context, resource, key, action string, extra logrus.Fields) {
	ctx, cancel := context.WithTimeout(parent, sideEffectTimeout)
	defer cancel()

	g, gctx := errgroup.WithContext(ctx)
	g.Go(guardStage("emit", resource, key, action, func() error { return emitEvent(gctx, resource, key, action, extra) }))
	g.Go(guardStage("cache", resource, key, action, func() error { return invalidateCache(gctx, resource, key) }))
	g.Go(guardStage("audit", resource, key, action, func() error { return recordAudit(gctx, resource, key, action, extra) }))

	if err := g.Wait(); err != nil {
		logrus.WithError(err).
			WithField("resource", resource).
			WithField("key", key).
			WithField("action", action).
			Warn("side effects partially failed")
	}
}

func guardStage(stage, resource, key, action string, fn func() error) func() error {
	return func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				logrus.WithField("panic", r).
					WithField("stage", stage).
					WithField("resource", resource).
					WithField("key", key).
					WithField("action", action).
					Error("side effect stage panicked; recovered")
				err = fmt.Errorf("side effect %s panicked: %v", stage, r)
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
	logrus.WithContext(ctx).
		WithField("resource", resource).
		WithField("key", key).
		Info("cache invalidated")
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
