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
	middleware.RegisterMetrics()

	e := gin.New()
	e.Use(middleware.ZapLoggerMiddleware())
	e.Use(middleware.ZapRecoveryMiddleware())
	e.Use(middleware.MetricsMiddleware())

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

	e.GET("/health/liveness", healthController.GetLiveness)
	e.GET("/health/readiness", healthController.GetReadiness)

	e.GET("/foos", fooController.GetAll)
	e.GET("foos/:id", fooController.GetByID)
	e.POST("/foos", fooController.Create)
	e.PUT("/foos/:id", fooController.Update)
	e.DELETE("/foos/:id", fooController.DeleteByID)

	return e
}
