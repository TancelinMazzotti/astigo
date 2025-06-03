package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type HealthController struct{}

func (c *HealthController) GetLiveness(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":     "UP",
		"message":    "Service is healthy",
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
		"start_time": StartAt.Format(time.RFC3339),
	})
}

func (c *HealthController) GetReadiness(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":     "UP",
		"message":    "Service is healthy",
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
		"start_time": StartAt.Format(time.RFC3339),
	})
}

func NewHealthController() *HealthController {
	c := &HealthController{}

	return c
}
