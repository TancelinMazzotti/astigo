package http

import (
	"astigo/internal/application/http/dto"
	"astigo/internal/domain/handler"
	"astigo/internal/domain/repository"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type FooController struct {
	svc handler.IFooHandler
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	foos, err := c.svc.GetAll(ctx, handler.PaginationInput{
		Offset: queryParams.Offset,
		Limit:  queryParams.Limit,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	results := make([]dto.FooReadResponse, len(foos))
	for i, foo := range foos {
		results[i] = dto.FooReadResponse{
			Id:    foo.Id,
			Label: foo.Label,
		}
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(pathParams.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	foo, err := c.svc.GetByID(ctx, id)
	if err != nil {

		if errors.As(err, &repository.ErrorNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := dto.FooReadResponse{
		Id:    foo.Id,
		Label: foo.Label,
	}
	ctx.JSON(http.StatusOK, result)
}

// Create @Summary Create a new foo
// @Description Create a new foo
// @Tags Foo
// @Accept JSON
// @Produce JSON
// @Param foo body dto.FooCreateRequest true "Foo"
// @Success 201
// @Router /foos [post]
func (c *FooController) Create(ctx *gin.Context) {
	var input dto.FooCreateRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	foo, err := c.svc.Create(ctx, handler.FooCreateInput{
		Label:  input.Label,
		Secret: input.Secret,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"id": foo.Id.String()})
}

// Update @Summary Update a foo
// @Description Update a foo
// @Tags Foo
// @Accept JSON
// @Produce JSON
// @Param foo body dto.FooUpdateRequest true "Foo"
// @Success 204
// @Router /foos [put]
func (c *FooController) Update(ctx *gin.Context) {
	var input dto.FooUpdateRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.svc.Update(ctx, handler.FooUpdateInput{
		Id:     input.Id,
		Label:  input.Label,
		Secret: input.Secret,
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

// DeleteByID @Summary Delete a foo
// @Description Delete a foo
// @Tags Foo
// @Accept JSON
// @Produce JSON
// @Param id path int true "Foo id"
// @Success 204
// @Router /foos/{id} [delete]
func (c *FooController) DeleteByID(ctx *gin.Context) {
	var pathParams dto.FooReadRequest

	if err := ctx.ShouldBindUri(&pathParams); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(pathParams.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.svc.DeleteByID(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func NewFooController(svc handler.IFooHandler) *FooController {
	c := &FooController{
		svc: svc,
	}

	return c
}
