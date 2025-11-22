package http

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/TancelinMazzotti/astigo/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/TancelinMazzotti/astigo/internal/application/http/middleware"
	"github.com/TancelinMazzotti/astigo/internal/domain/model"
	"github.com/TancelinMazzotti/astigo/internal/domain/service"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
)

var StartAt time.Time

type Config struct {
	Port     string `mapstructure:"port"`
	Mode     string `mapstructure:"mode"`
	Issuer   string `mapstructure:"issuer"`
	ClientID string `mapstructure:"client_id"`
}

func NewGin(
	config Config,
	logger *zap.Logger,
	authHandler service.IAuthService,
	healthController *HealthController,
	fooController *FooController,
) *gin.Engine {

	middleware.RegisterMetrics()
	gin.SetMode(config.Mode)
	authMiddleware := middleware.NewAuthMiddleware(authHandler)

	e := gin.New()
	e.Use(otelgin.Middleware("astigo"))
	e.Use(middleware.ZapLoggerMiddleware(logger))
	e.Use(middleware.ZapRecoveryMiddleware(logger))
	e.Use(middleware.MetricsMiddleware())
	e.Use(middleware.CorsMiddleware())

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

	e.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	e.GET("/metrics", gin.WrapH(promhttp.Handler()))

	e.GET("/health/liveness", healthController.GetLiveness)
	e.GET("/health/readiness", healthController.GetReadiness)

	e.GET("/foos", fooController.GetAll)
	e.GET("foos/:id", fooController.GetByID)
	e.POST("/foos", fooController.Create)
	e.PUT("/foos/:id", fooController.Update)
	e.PATCH("/foos/:id", fooController.Patch)
	e.DELETE("/foos/:id", fooController.DeleteByID)

	e.GET("/private", authMiddleware.Middleware, func(c *gin.Context) {
		claimsCtx, _ := c.Get("claims")
		claims, ok := claimsCtx.(*model.Claims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("invalid claims type")})
		}

		c.JSON(http.StatusOK, gin.H{
			"claims": claims,
		})
	})

	return e
}
