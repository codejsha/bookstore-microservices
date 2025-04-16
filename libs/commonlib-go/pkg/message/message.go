package message

import (
	"context"
	"time"
)

type EventMessage interface {
	Id() string
	Type() string
	Metadata() map[string]interface{}
	CreatedAt() time.Time
}

type EventHandler interface {
	HandleMessage(ctx context.Context, msg EventMessage) error
}

type EventPublisher interface {
	Publish(ctx context.Context, topicName string, msg EventMessage) error
}

type EventSubscriber interface {
	Subscribe(topicName string, handler EventHandler) error
	Unsubscribe() error
}

type BaseEvent struct {
	Id        string                 `json:"id" binding:"required,uuid"`
	Type      string                 `json:"type" binding:"required"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt time.Time              `json:"created_at" binding:"required"`
	EventMessage
}
