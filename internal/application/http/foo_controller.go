package http

import (
	"errors"
	"net/http"

	"github.com/TancelinMazzotti/astigo/internal/application/http/dto"
	"github.com/TancelinMazzotti/astigo/internal/domain/port"
	"github.com/TancelinMazzotti/astigo/internal/domain/port/in/data"
	"github.com/TancelinMazzotti/astigo/internal/domain/port/in/service"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var _ IFooController = (*FooController)(nil)

// IFooController defines an interface for managing Foo entity operations through HTTP handlers.
// GetAll retrieves all Foo entities.
// GetByID retrieves a Foo entity by its unique identifier.
// Create handles the creation of a new Foo entity.
// Update modifies an existing Foo entity.
// DeleteByID deletes a Foo entity by its unique identifier.
type IFooController interface {
	GetAll(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Patch(ctx *gin.Context)
	DeleteByID(ctx *gin.Context)
}

// FooController manages the HTTP request handling for operations related to Foo entities.
type FooController struct {
	svc service.IFooService
}

// GetAll @Summary Get all foo
// @Description Get all foos
// @Tags Foo
// @Accept JSON
// @Produce JSON
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Success 200 {array} dto.FooReadResponse
// @Router /foos [get]
func (c *FooController) GetAll(ctx *gin.Context) {
	tracer := otel.Tracer("FooController")
	spanCtx, span := tracer.Start(ctx.Request.Context(), "FooController.GetAll")
	defer span.End()

	var queryParams dto.ListRequest

	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to validate query params")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate query params"})
		return
	}

	foos, err := c.svc.GetAll(spanCtx, data.FooReadListInput{
		Offset: queryParams.Offset,
		Limit:  queryParams.Limit,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get all foos")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get all foos"})
		return
	}

	results := make([]*dto.FooReadResponse, len(foos))
	for i, foo := range foos {
		results[i] = dto.NewFooReadResponse(foo)
	}

	span.SetStatus(codes.Ok, "")
	span.SetAttributes(attribute.Int("response.count", len(results)))
	ctx.JSON(http.StatusOK, results)
}

// GetByID @Summary Get foo by id
// @Description Get foo by id
// @Tags Foo
// @Accept JSON
// @Produce JSON
// @Param id path uuid true "Foo id"
// @Success 200 {object} dto.FooReadResponse
// @Router /foos/{id} [get]
func (c *FooController) GetByID(ctx *gin.Context) {
	tracer := otel.Tracer("FooController")
	spanCtx, span := tracer.Start(ctx.Request.Context(), "FooController.GetById")
	defer span.End()

	var pathParams dto.FooReadRequest

	if err := ctx.ShouldBindUri(&pathParams); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to validate path params")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate path params"})
		return
	}

	id, err := uuid.Parse(pathParams.Id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to parse id to uuid")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse id to uuid"})
		return
	}
	span.SetAttributes(attribute.String("foo.id", id.String()))

	foo, err := c.svc.GetByID(spanCtx, id)
	if err != nil {
		span.RecordError(err)

		if errors.As(err, &port.ErrorNotFound) {
			span.SetStatus(codes.Error, "foo not found")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "foo not found"})
			return
		}
		span.SetStatus(codes.Error, "failed to get foo by id")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get foo by id"})
		return
	}

	result := dto.NewFooReadResponse(foo)

	span.SetStatus(codes.Ok, "")
	span.SetAttributes(
		attribute.String("foo.label", foo.Label),
		attribute.Int("foo.value", foo.Value),
		attribute.Float64("foo.weight", float64(foo.Weight)),
	)

	ctx.JSON(http.StatusOK, result)
}

// Create @Summary Create a new foo
// @Description Create a new foo
// @Tags Foo
// @Accept JSON
// @Produce JSON
// @Param foo body dto.FooCreateBody true "Foo"
// @Success 201 {object} dto.FooCreateResponse
// @Router /foos [post]
func (c *FooController) Create(ctx *gin.Context) {
	tracer := otel.Tracer("FooController")
	spanCtx, span := tracer.Start(ctx.Request.Context(), "FooController.Create")
	defer span.End()

	var input dto.FooCreateBody
	if err := ctx.ShouldBindJSON(&input); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to validate request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate request body"})
		return
	}

	span.SetAttributes(
		attribute.String("foo.label", input.Label),
		attribute.Int("foo.value", input.Value),
		attribute.Float64("foo.weight", float64(input.Weight)),
	)

	foo, err := c.svc.Create(spanCtx, data.FooCreateInput{
		Label:  input.Label,
		Secret: input.Secret,
		Value:  input.Value,
		Weight: input.Weight,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create foo")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create foo"})
		return
	}

	result := &dto.FooCreateResponse{
		Id: foo.Id,
	}

	span.SetStatus(codes.Ok, "")
	span.SetAttributes(attribute.String("foo.id", foo.Id.String()))
	ctx.JSON(http.StatusCreated, result)
}

// Update @Summary Update a foo
// @Description Update a foo
// @Tags Foo
// @Accept JSON
// @Produce JSON
// @Param foo body dto.FooUpdateBody true "Foo"
// @Success 204
// @Router /foos/{id} [put]
func (c *FooController) Update(ctx *gin.Context) {
	tracer := otel.Tracer("FooController")
	spanCtx, span := tracer.Start(ctx.Request.Context(), "FooController.Update")
	defer span.End()

	var pathParams dto.FooUpdateRequest
	var body dto.FooUpdateBody
	if err := ctx.ShouldBindUri(&pathParams); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to validate path params")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate path params"})
		return
	}

	id, err := uuid.Parse(pathParams.Id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to parse id to uuid")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse id to uuid"})
		return
	}
	span.SetAttributes(attribute.String("foo.id", id.String()))

	if err := ctx.ShouldBindJSON(&body); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to validate request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate request body"})
		return
	}

	span.SetAttributes(
		attribute.String("foo.label", body.Label),
		attribute.Int("foo.value", body.Value),
		attribute.Float64("foo.weight", float64(body.Weight)),
	)

	if err := c.svc.Update(spanCtx, &data.FooUpdateInput{
		Id:     id,
		Label:  body.Label,
		Secret: body.Secret,
		Value:  body.Value,
		Weight: body.Weight,
	}); err != nil {
		span.RecordError(err)
		if errors.As(err, &port.ErrorNotFound) {
			span.SetStatus(codes.Error, "foo not found")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "foo not found"})
			return
		}
		span.SetStatus(codes.Error, "failed to update foo")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update foo"})
		return
	}

	span.SetStatus(codes.Ok, "")
	ctx.Status(http.StatusNoContent)
}

// Patch @Summary Patch a foo
// @Description Patch a foo
// @Tags Foo
// @Accept JSON
// @Produce JSON
// @Param foo body dto.FooPatchBody true "Foo"
// @Success 204
// @Router /foos/{id} [patch]
func (c *FooController) Patch(ctx *gin.Context) {
	tracer := otel.Tracer("FooController")
	spanCtx, span := tracer.Start(ctx.Request.Context(), "FooController.Patch")
	defer span.End()

	var pathParams dto.FooPatchRequest
	var body dto.FooPatchBody
	if err := ctx.ShouldBindUri(&pathParams); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to validate path params")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate path params"})
		return
	}

	id, err := uuid.Parse(pathParams.Id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to parse id to uuid")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse id to uuid"})
		return
	}
	span.SetAttributes(attribute.String("foo.id", id.String()))

	if err := ctx.ShouldBindJSON(&body); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to validate request body")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate request body"})
		return
	}

	attrs := make([]attribute.KeyValue, 0)
	if body.Label != nil {
		attrs = append(attrs, attribute.String("foo.label", *body.Label))
	}
	if body.Value != nil {
		attrs = append(attrs, attribute.Int("foo.value", *body.Value))
	}
	if body.Weight != nil {
		attrs = append(attrs, attribute.Float64("foo.weight", float64(*body.Weight)))
	}
	span.SetAttributes(attrs...)

	var input data.FooPatchInput
	input.Id = id
	if body.Label != nil {
		input.Label.Set = true
		input.Label.Value = *body.Label
	}
	if body.Secret != nil {
		input.Secret.Set = true
		input.Secret.Value = *body.Secret
	}
	if body.Value != nil {
		input.Value.Set = true
		input.Value.Value = *body.Value
	}
	if body.Weight != nil {
		input.Weight.Set = true
		input.Weight.Value = *body.Weight
	}

	if err := c.svc.Update(spanCtx, &input); err != nil {
		span.RecordError(err)
		if errors.As(err, &port.ErrorNotFound) {
			span.SetStatus(codes.Error, "foo not found")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "foo not found"})
			return
		}
		span.SetStatus(codes.Error, "failed to update foo")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update foo"})
		return
	}

	span.SetStatus(codes.Ok, "")
	ctx.Status(http.StatusNoContent)
}

// DeleteByID @Summary Delete a foo
// @Description Delete a foo
// @Tags Foo
// @Accept JSON
// @Produce JSON
// @Param id path uuid true "Foo id"
// @Success 204
// @Router /foos/{id} [delete]
func (c *FooController) DeleteByID(ctx *gin.Context) {
	tracer := otel.Tracer("FooController")
	spanCtx, span := tracer.Start(ctx.Request.Context(), "FooController.DeleteByID")
	defer span.End()

	var pathParams dto.FooDeleteRequest

	if err := ctx.ShouldBindUri(&pathParams); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to validate path params")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate path params"})
		return
	}

	id, err := uuid.Parse(pathParams.Id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to validate path params")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse id to uuid"})
		return
	}
	span.SetAttributes(attribute.String("foo.id", id.String()))

	if err := c.svc.DeleteByID(spanCtx, id); err != nil {
		span.RecordError(err)
		if errors.As(err, &port.ErrorNotFound) {
			span.SetStatus(codes.Error, "foo not found")
			ctx.JSON(http.StatusNotFound, gin.H{"error": "foo not found"})
			return
		}
		span.SetStatus(codes.Error, "failed to delete foo")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete foo"})
		return
	}

	span.SetStatus(codes.Ok, "")
	ctx.Status(http.StatusNoContent)
}

// NewFooController initializes a new FooController with the provided IFooService dependency.
func NewFooController(svc service.IFooService) *FooController {
	c := &FooController{
		svc: svc,
	}

	return c
}
