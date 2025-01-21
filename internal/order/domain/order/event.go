package order

import "context"

type DomainEvent struct {
	Dest string
	Data any
}

type EventPublisher interface {
	Publish(ctx context.Context, event DomainEvent) error
	Broadcast(ctx context.Context, event DomainEvent) error
}
