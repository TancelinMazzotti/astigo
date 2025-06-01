package http

import (
	"astigo/internal/application/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

type RouterRegistrar interface {
	RegisterRoutes(router *gin.Engine)
}

type GinConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

func NewGin() *gin.Engine {
	r := gin.New()
	r.Use(middleware.ZapLoggerMiddleware())
	r.Use(middleware.ZapRecoveryMiddleware())

	setupSystemRoutes(r)

	return r
}

func setupSystemRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":        "Astigo API",
			"version":     "1.0.0",
			"description": "Welcome to the Astigo API. This is a RESTful service.",
		})
	})

	router.GET("/health/liveness", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":     "UP",
			"message":    "Service is healthy",
			"timestamp":  time.Now().UTC().Format(time.RFC3339),
			"start_time": time.Now().UTC().Format(time.RFC3339),
		})
	})

	router.GET("/health/readiness", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":     "UP",
			"message":    "Service is healthy",
			"timestamp":  time.Now().UTC().Format(time.RFC3339),
			"start_time": time.Now().UTC().Format(time.RFC3339),
		})
	})

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

}
