package middleware

import (
	"astigo/internal/domain/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// AuthMiddleware handles authentication by validating JWT tokens in incoming requests.
// It uses a provided implementation of IAuthService for token verification and claims extraction.
type AuthMiddleware struct {
	handler service.IAuthService
}

// Middleware is a Gin middleware function that validates JWT Authorization headers for protected routes.
func (m *AuthMiddleware) Middleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header"})
		return
	}

	token := parts[1]

	idToken, err := m.handler.VerifyToken(c, token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, err := m.handler.GetClaims(idToken)
	if err != nil {

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse claims"})
		return
	}

	c.Set("claims", claims)

	c.Next()

}

// NewAuthMiddleware creates and returns an instance of AuthMiddleware using the provided IAuthService for authentication.
func NewAuthMiddleware(authHandler service.IAuthService) *AuthMiddleware {
	return &AuthMiddleware{
		handler: authHandler,
	}
}
