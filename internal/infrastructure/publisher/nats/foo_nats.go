package nats

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/TancelinMazzotti/astigo/internal/domain/model"
	"github.com/TancelinMazzotti/astigo/internal/domain/port/out/publisher"
	"github.com/TancelinMazzotti/astigo/internal/infrastructure/publisher/nats/message"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

const (
	fooCreatedSubject = "foo.created"
	fooUpdatedSubject = "foo.updated"
	fooDeletedSubject = "foo.deleted"
)

var (
	_ publisher.IFooPublisher = (*FooNats)(nil)
)

// FooNats wraps a NATS connection and implements the IFooPublisher interface for publishing Foo-related messages.
type FooNats struct {
	conn *nats.Conn
}

// PublishCreated publishes a "foo.created" message to the NATS server using the provided Foo data.
func (n *FooNats) PublishCreated(ctx context.Context, foo *model.Foo) error {
	tracer := otel.Tracer("FooNats")
	_, span := tracer.Start(ctx, "FooNats.PublishCreated")
	defer span.End()

	span.SetAttributes(
		attribute.String("foo.id", foo.Id.String()),
		attribute.String("foo.label", foo.Label),
		attribute.Int("foo.value", foo.Value),
		attribute.Float64("foo.weight", float64(foo.Weight)),
		attribute.String("nats.subject", fooCreatedSubject),
	)

	msg := message.NewFooMessage(foo)
	data, err := json.Marshal(msg)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to serialize foo")
		return fmt.Errorf("failed to serialize Foo: %w", err)
	}

	span.SetAttributes(attribute.Int("message.size", len(data)))

	if err := n.conn.Publish(fooCreatedSubject, data); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to publish message")
		return fmt.Errorf("failed to publish to NATS: %w", err)
	}

	span.SetStatus(codes.Ok, "")

	return nil
}

func (n *FooNats) PublishUpdated(ctx context.Context, foo *model.Foo) error {
	tracer := otel.Tracer("FooNats")
	_, span := tracer.Start(ctx, "FooNats.PublishUpdated")
	defer span.End()

	span.SetAttributes(
		attribute.String("foo.id", foo.Id.String()),
		attribute.String("foo.label", foo.Label),
		attribute.Int("foo.value", foo.Value),
		attribute.Float64("foo.weight", float64(foo.Weight)),
		attribute.String("nats.subject", fooUpdatedSubject),
	)

	msg := message.NewFooMessage(foo)
	data, err := json.Marshal(msg)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to serialize foo")
		return fmt.Errorf("failed to serialize Foo: %w", err)
	}

	span.SetAttributes(attribute.Int("message.size", len(data)))

	if err := n.conn.Publish(fooUpdatedSubject, data); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to publish message")
		return fmt.Errorf("failed to publish to NATS: %w", err)
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

func (n *FooNats) PublishDeleted(ctx context.Context, id uuid.UUID) error {
	tracer := otel.Tracer("FooNats")
	_, span := tracer.Start(ctx, "FooNats.PublishDeleted")
	defer span.End()

	span.SetAttributes(
		attribute.String("foo.id", id.String()),
		attribute.String("nats.subject", fooDeletedSubject),
	)

	data, err := json.Marshal(map[string]string{"id": id.String()})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to serialize id")
		return fmt.Errorf("failed to serialize ID: %w", err)
	}

	span.SetAttributes(attribute.Int("message.size", len(data)))

	if err := n.conn.Publish(fooDeletedSubject, data); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to publish message")
		return fmt.Errorf("failed to publish to NATS: %w", err)
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

func NewFooNats(conn *nats.Conn) *FooNats {
	return &FooNats{conn: conn}
}
