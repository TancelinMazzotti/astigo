package http

import (
	"astigo/internal/domain/handler"
	"astigo/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FooController struct {
	svc handler.IFooHandler
}

func (c *FooController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/:id", c.GetByID)
	r.POST("/", c.Create)
}

func (c *FooController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.svc.Get(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *FooController) Create(ctx *gin.Context) {
	var input model.Foo
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.svc.Register(ctx, input); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusCreated)
}

func NewFooController(svc handler.IFooHandler) *FooController {
	return &FooController{svc: svc}
}
