package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/TancelinMazzotti/astigo/internal/domain/model"
	"github.com/TancelinMazzotti/astigo/internal/domain/port/in/data"
	"github.com/TancelinMazzotti/astigo/internal/domain/port/in/service"
	"github.com/TancelinMazzotti/astigo/internal/domain/port/out/cache"
	"github.com/TancelinMazzotti/astigo/internal/domain/port/out/messaging"
	"github.com/TancelinMazzotti/astigo/internal/domain/port/out/repository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const FooCacheExpiration = time.Minute * 15

var (
	_ service.IFooService = (*FooService)(nil)
)

// FooService provides business logic around Foo entities, integrating data access, caching, and messaging capabilities.
type FooService struct {
	logger    *zap.Logger
	repo      repository.IFooRepository
	cache     cache.IFooCache
	messaging messaging.IFooMessaging
}

// GetAll retrieves a list of Foo entities based on the provided input criteria and returns an error if retrieval fails.
func (s *FooService) GetAll(ctx context.Context, input data.FooReadListInput) ([]*model.Foo, error) {
	tracer := otel.Tracer("FooService")
	ctx, span := tracer.Start(ctx, "FooService.GetAll")
	defer span.End()

	span.SetAttributes(
		attribute.Int("offset", input.Offset),
		attribute.Int("limit", input.Limit),
	)

	foos, err := s.repo.FindAll(ctx, input)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to find all foo")

		s.logger.Debug("fail to find all foo", zap.Error(err))
		return nil, fmt.Errorf("fail to find all foo: %w", err)
	}

	span.SetStatus(codes.Ok, "")
	span.SetAttributes(attribute.Int("result.count", len(foos)))
	return foos, nil
}

// GetByID retrieves a Foo entity by its ID, using a cache-first approach and falling back to the repository if needed.
func (s *FooService) GetByID(ctx context.Context, id uuid.UUID) (*model.Foo, error) {
	tracer := otel.Tracer("FooService")
	ctx, span := tracer.Start(ctx, "FooService.GetByID")
	defer span.End()

	span.SetAttributes(attribute.String("id", id.String()))

	foo, err := s.cache.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Bool("cache.get.error", true))
		s.logger.Debug("fail to find foo by id from cache", zap.Error(err))
	}

	if foo == nil {
		span.SetAttributes(attribute.Bool("cache.miss", true))

		foo, err = s.repo.FindByID(ctx, id)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to find foo by id")
			s.logger.Debug("fail to find foo by id", zap.Error(err))
			return nil, fmt.Errorf("fail to find foo by id: %w", err)
		}

		if err := s.cache.Set(ctx, foo, FooCacheExpiration); err != nil {
			span.RecordError(err)
			span.SetAttributes(attribute.Bool("cache.set.error", true))
			s.logger.Warn("fail to create foo in cache: %w", zap.Error(err))
		}
	} else {
		span.SetAttributes(attribute.Bool("cache.hit", true))
	}

	return foo, nil
}

// Create creates a new Foo entity, stores it in the repository, and updates related cache and messaging.
func (s *FooService) Create(ctx context.Context, input data.FooCreateInput) (*model.Foo, error) {
	tracer := otel.Tracer("FooService")
	ctx, span := tracer.Start(ctx, "FooService.Create")
	defer span.End()

	foo := &model.Foo{
		Id:     uuid.New(),
		Label:  input.Label,
		Secret: input.Secret,
		Value:  input.Value,
		Weight: input.Weight,
	}

	span.SetAttributes(
		attribute.String("foo.id", foo.Id.String()),
		attribute.String("foo.label", foo.Label),
		attribute.Int("foo.value", foo.Value),
		attribute.Float64("foo.weight", float64(foo.Weight)),
	)

	var validate = validator.New()
	if err := validate.Struct(foo); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid input")
		s.logger.Debug("invalid input", zap.Error(err))
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	if err := s.repo.Create(ctx, foo); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "fail to create foo")
		s.logger.Debug("fail to create foo", zap.Error(err))
		return nil, fmt.Errorf("fail to create foo: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := s.cache.Set(ctx, foo, FooCacheExpiration); err != nil {
			span.RecordError(err)
			span.SetAttributes(attribute.Bool("cache.set.error", true))
			s.logger.Warn("fail to create foo in cache", zap.Error(err))
		}
	}()

	var errMessaging error
	go func() {
		defer wg.Done()
		if err := s.messaging.PublishFooCreated(ctx, foo); err != nil {
			span.RecordError(err)
			span.SetAttributes(attribute.Bool("messaging.publish.error", true))
			errMessaging = err
		}

	}()

	wg.Wait()

	if errMessaging != nil {
		s.logger.Debug("fail to publish foo created", zap.Error(errMessaging))
		return nil, fmt.Errorf("fail to publish foo created: %w", errMessaging)
	}

	span.SetStatus(codes.Ok, "")
	return foo, nil
}

// Update applies partial updates to an existing Foo entity based on the provided input and propagates changes across systems.
// It retrieves the entity by ID, merges changes, updates the repository, cache, and publishes an event.
func (s *FooService) Update(ctx context.Context, input data.IFooUpdateMerger) error {
	tracer := otel.Tracer("FooService")
	ctx, span := tracer.Start(ctx, "FooService.Update")
	defer span.End()

	span.SetAttributes(attribute.String("foo.id", input.GetID().String()))

	foo, err := s.repo.FindByID(ctx, input.GetID())
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "fail to find foo by id")
		s.logger.Debug("fail to find foo by id", zap.Error(err))
		return fmt.Errorf("fail to get foo by id: %w", err)
	}

	oldValues := map[string]interface{}{
		"label":  foo.Label,
		"value":  foo.Value,
		"weight": foo.Weight,
	}

	if err := input.Merge(foo); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "fail to merge input")
		s.logger.Debug("fail to merge input", zap.Error(err))
		return fmt.Errorf("fail to merge input: %w", err)
	}

	span.SetAttributes(
		attribute.String("update.label.old", oldValues["label"].(string)),
		attribute.String("update.label.new", foo.Label),
		attribute.Int("update.value.old", oldValues["value"].(int)),
		attribute.Int("update.value.new", foo.Value),
		attribute.Float64("update.weight.old", float64(oldValues["weight"].(float32))),
		attribute.Float64("update.weight.new", float64(foo.Weight)),
	)

	var validate = validator.New()
	if err := validate.Struct(foo); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid input")
		s.logger.Debug("invalid input", zap.Error(err))
		return fmt.Errorf("invalid input: %w", err)
	}

	if err := s.repo.Update(ctx, foo); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "fail to update foo")
		s.logger.Debug("fail to update foo", zap.Error(err))
		return fmt.Errorf("fail to update foo: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := s.cache.Set(ctx, foo, FooCacheExpiration); err != nil {
			span.RecordError(err)
			span.SetAttributes(attribute.Bool("cache.set.error", true))
			s.logger.Warn("fail to update foo in cache", zap.Error(err))
		}
	}()

	var errMessaging error
	go func() {
		defer wg.Done()
		if err := s.messaging.PublishFooUpdated(ctx, foo); err != nil {
			span.RecordError(err)
			span.SetAttributes(attribute.Bool("messaging.publish.error", true))
			errMessaging = err
		}

	}()

	wg.Wait()

	if errMessaging != nil {
		s.logger.Debug("fail to publish foo updated", zap.Error(errMessaging))
		return fmt.Errorf("fail to publish foo updated: %w", errMessaging)
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

// DeleteByID removes a Foo entity by its ID, updates the cache, and publishes a deletion event. Returns an error if any step fails.
func (s *FooService) DeleteByID(ctx context.Context, id uuid.UUID) error {
	tracer := otel.Tracer("FooService")
	ctx, span := tracer.Start(ctx, "FooService.DeleteByID")
	defer span.End()

	span.SetAttributes(attribute.String("foo.id", id.String()))

	if err := s.repo.DeleteByID(ctx, id); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "fail to delete foo")
		s.logger.Debug("fail to delete foo by id", zap.Error(err))
		return fmt.Errorf("fail to delete foo by id: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := s.cache.DeleteByID(ctx, id); err != nil {
			span.RecordError(err)
			span.SetAttributes(attribute.Bool("cache.delete.error", true))
			s.logger.Warn("fail to delete foo by id from cache", zap.Error(err))
		}
	}()

	var errMessaging error
	go func() {
		defer wg.Done()
		if err := s.messaging.PublishFooDeleted(ctx, id); err != nil {
			span.RecordError(err)
			span.SetAttributes(attribute.Bool("messaging.publish.error", true))
			errMessaging = err
		}
	}()

	wg.Wait()

	if errMessaging != nil {
		s.logger.Debug("fail to publish foo deleted", zap.Error(errMessaging))
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

// NewFooService initializes a new instance of FooService with the provided logger, repository, cache, and messaging dependencies.
func NewFooService(logger *zap.Logger, repo repository.IFooRepository, cache cache.IFooCache, messaging messaging.IFooMessaging) *FooService {
	return &FooService{
		logger:    logger,
		repo:      repo,
		cache:     cache,
		messaging: messaging,
	}
}
