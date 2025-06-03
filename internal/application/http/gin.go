package http

import (
	"astigo/internal/application/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var StartAt time.Time

type GinConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

func NewGin(healthController *HealthController, fooController *FooController) *gin.Engine {
	e := gin.New()
	e.Use(middleware.ZapLoggerMiddleware())
	e.Use(middleware.ZapRecoveryMiddleware())

	e.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":        "Astigo API",
			"version":     "1.0.0",
			"description": "Welcome to the Astigo API. This is a RESTful service.",
			"author":      "Tancelin Mazzotti",
			"github":      "https://github.com/TancelinMazzotti/astigo",
			"license":     "MIT",
			"docs":        "https://github.com/TancelinMazzotti/astigo/blob/main/README.md",
		})
	})

	e.GET("/metrics", gin.WrapH(promhttp.Handler()))

	{
		health := e.Group("/health")
		health.GET("/liveness", healthController.GetLiveness)
		health.GET("/readiness", healthController.GetReadiness)
	}

	{
		foos := e.Group("/foos")
		foos.GET("", fooController.GetAll)
		foos.GET("/:id", fooController.GetByID)
		foos.POST("", fooController.Create)
		foos.PUT("", fooController.Update)
		foos.DELETE("/:id", fooController.DeleteByID)

	}

	return e
}
