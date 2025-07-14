package http

import (
	"astigo/internal/application/http/dto"
	"astigo/internal/domain/adapter/data"
	"astigo/internal/domain/adapter/repository"
	"astigo/internal/domain/service"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
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
	var queryParams dto.ListRequest

	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate query params"})
		return
	}

	foos, err := c.svc.GetAll(ctx, data.FooReadListInput{
		Offset: queryParams.Offset,
		Limit:  queryParams.Limit,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get all foos"})
		return
	}

	results := make([]*dto.FooReadResponse, len(foos))
	for i, foo := range foos {
		results[i] = dto.NewFooReadResponse(foo)
	}

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
	var pathParams dto.FooReadRequest

	if err := ctx.ShouldBindUri(&pathParams); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate path params"})
		return
	}

	id, err := uuid.Parse(pathParams.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse id to uuid"})
		return
	}

	foo, err := c.svc.GetByID(ctx, id)
	if err != nil {
		if errors.As(err, &repository.ErrorNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "foo not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get foo by id"})
		return
	}

	result := dto.NewFooReadResponse(foo)
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
	var input dto.FooCreateBody
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate request body"})
		return
	}

	foo, err := c.svc.Create(ctx, data.FooCreateInput{
		Label:  input.Label,
		Secret: input.Secret,
		Value:  input.Value,
		Weight: input.Weight,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create foo"})
		return
	}

	result := &dto.FooCreateResponse{
		Id: foo.Id,
	}
	ctx.JSON(http.StatusCreated, result)
}

// Update @Summary Update a foo
// @Description Update a foo
// @Tags Foo
// @Accept JSON
// @Produce JSON
// @Param foo body dto.FooUpdateBody true "Foo"
// @Success 204
// @Router /foos [put]
func (c *FooController) Update(ctx *gin.Context) {
	var pathParams dto.FooUpdateRequest
	var body dto.FooUpdateBody
	if err := ctx.ShouldBindUri(&pathParams); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate path params"})
		return
	}

	id, err := uuid.Parse(pathParams.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse id to uuid"})
		return
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate request body"})
		return
	}

	if err := c.svc.Update(ctx, data.FooUpdateInput{
		Id:     id,
		Label:  body.Label,
		Secret: body.Secret,
		Value:  body.Value,
		Weight: body.Weight,
	}); err != nil {
		if errors.As(err, &repository.ErrorNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "foo not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update foo"})
		return
	}
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
	var pathParams dto.FooDeleteRequest

	if err := ctx.ShouldBindUri(&pathParams); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to validate path params"})
		return
	}

	id, err := uuid.Parse(pathParams.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse id to uuid"})
		return
	}

	if err := c.svc.DeleteByID(ctx, id); err != nil {
		if errors.As(err, &repository.ErrorNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "foo not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete foo"})
		return
	}
	ctx.Status(http.StatusNoContent)
}

// NewFooController initializes a new FooController with the provided IFooService dependency.
func NewFooController(svc service.IFooService) *FooController {
	c := &FooController{
		svc: svc,
	}

	return c
}
