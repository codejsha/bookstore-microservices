package restcontroller

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/codejsha/shared-library-go/pkg/message"
)

const sideEffectTimeout = 5 * time.Second

var eventPublisher atomic.Pointer[message.EventPublisher]

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
	g.Go(func() error { return emitEvent(gctx, resource, key, action, extra) })
	g.Go(func() error { return invalidateCache(gctx, resource, key) })
	g.Go(func() error { return recordAudit(gctx, resource, key, action, extra) })

	if err := g.Wait(); err != nil {
		logrus.WithError(err).
			WithField("resource", resource).
			WithField("key", key).
			WithField("action", action).
			Warn("side effects partially failed")
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
