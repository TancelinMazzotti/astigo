package http

import (
	"astigo/internal/domain/handler"
	"astigo/pkg/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FooController struct {
	svc handler.IFooHandler
}

func (c *FooController) GetAll(ctx *gin.Context) {
	var queryParams dto.PaginationRequestDto

	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := c.svc.GetAll(ctx, queryParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, results)
}

func (c *FooController) GetByID(ctx *gin.Context) {
	var pathParams dto.FooRequestReadDto

	if err := ctx.ShouldBindUri(&pathParams); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.svc.GetByID(ctx, pathParams.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *FooController) Create(ctx *gin.Context) {
	var input dto.FooRequestCreateDto
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.svc.Create(ctx, input); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusCreated)
}

func (c *FooController) Update(ctx *gin.Context) {
	var input dto.FooRequestUpdateDto
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.svc.Update(ctx, input); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (c *FooController) DeleteByID(ctx *gin.Context) {
	var pathParams dto.FooRequestReadDto

	if err := ctx.ShouldBindUri(&pathParams); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.svc.DeleteByID(ctx, pathParams.Id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func NewFooController(engine *gin.Engine, svc handler.IFooHandler) *FooController {
	c := &FooController{
		svc: svc,
	}

	engine.GET("/foos", c.GetAll)
	engine.GET("/foos/:id", c.GetByID)
	engine.POST("/foos", c.Create)
	engine.PUT("/foos", c.Update)
	engine.DELETE("/foos/:id", c.DeleteByID)

	return c
}
