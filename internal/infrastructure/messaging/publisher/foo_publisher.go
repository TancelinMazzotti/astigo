package publisher

import (
	"context"
	"fmt"

	"github.com/TancelinMazzotti/astigo/internal/domain/contract/messaging"
	"github.com/TancelinMazzotti/astigo/internal/domain/model"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"golang.org/x/sync/errgroup"
)

var _ messaging.IFooMessaging = (*FooPublisher)(nil)

type FooPublisher struct {
	Subscribers []messaging.IFooMessaging
}

// Subscribe to publisher
func (p *FooPublisher) Subscribe(subscriber messaging.IFooMessaging) {
	p.Subscribers = append(p.Subscribers, subscriber)
}

// Unsubscribe to publisher
func (p *FooPublisher) Unsubscribe(subscriber messaging.IFooMessaging) {
	for i := len(p.Subscribers) - 1; i >= 0; i-- {
		if p.Subscribers[i] == subscriber {
			p.Subscribers = append(p.Subscribers[:i], p.Subscribers[i+1:]...)
		}
	}
}

func (p *FooPublisher) PublishFooCreated(ctx context.Context, foo *model.Foo) error {
	tracer := otel.Tracer("FooPublisher")
	_, span := tracer.Start(ctx, "FooPublisher.PublishFooCreated")
	defer span.End()

	g, ctx := errgroup.WithContext(ctx)

	for _, subscriber := range p.Subscribers {
		g.Go(func() error {
			return subscriber.PublishFooCreated(ctx, foo)
		})
	}

	if err := g.Wait(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to publish foo created")
		return fmt.Errorf("failed to publish foo created: %w", err)
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

func (p *FooPublisher) PublishFooUpdated(ctx context.Context, foo *model.Foo) error {
	tracer := otel.Tracer("FooPublisher")
	_, span := tracer.Start(ctx, "FooPublisher.PublishFooCreated")
	defer span.End()

	g, ctx := errgroup.WithContext(ctx)

	for _, subscriber := range p.Subscribers {
		g.Go(func() error {
			return subscriber.PublishFooUpdated(ctx, foo)
		})
	}

	if err := g.Wait(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to publish foo updated")
		return fmt.Errorf("failed to publish foo updated: %w", err)
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

func (p *FooPublisher) PublishFooDeleted(ctx context.Context, id uuid.UUID) error {
	tracer := otel.Tracer("FooPublisher")
	_, span := tracer.Start(ctx, "FooPublisher.PublishFooCreated")
	defer span.End()

	g, ctx := errgroup.WithContext(ctx)

	for _, subscriber := range p.Subscribers {
		g.Go(func() error {
			return subscriber.PublishFooDeleted(ctx, id)
		})
	}

	if err := g.Wait(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to publish foo deleted")
		return fmt.Errorf("failed to publish foo deleted: %w", err)
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

func NewFooPublisher() *FooPublisher {
	return &FooPublisher{Subscribers: make([]messaging.IFooMessaging, 0)}
}
