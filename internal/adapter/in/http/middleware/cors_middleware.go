package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type CorsConfig struct {
	AllowAllOrigins           bool          `mapstructure:"allow_all_origins"`
	AllowOrigins              []string      `mapstructure:"allow_origins"`
	AllowMethods              []string      `mapstructure:"allow_methods"`
	AllowPrivateNetwork       bool          `mapstructure:"allow_private_network"`
	AllowHeaders              []string      `mapstructure:"allow_headers"`
	AllowCredentials          bool          `mapstructure:"allow_credentials"`
	ExposeHeaders             []string      `mapstructure:"expose_headers"`
	MaxAge                    time.Duration `mapstructure:"max_age"`
	AllowWildcard             bool          `mapstructure:"allow_wildcard"`
	AllowBrowserExtensions    bool          `mapstructure:"allow_browser_extensions"`
	CustomSchemas             []string      `mapstructure:"custom_schemas"`
	AllowWebSockets           bool          `mapstructure:"allow_websockets"`
	AllowFiles                bool          `mapstructure:"allow_files"`
	OptionsResponseStatusCode int           `mapstructure:"options_response_status_code"`
}

func CorsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
