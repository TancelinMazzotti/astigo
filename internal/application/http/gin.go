package http

import (
	"astigo/internal/application/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type GinConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

func NewGin() *gin.Engine {
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

	return e
}
