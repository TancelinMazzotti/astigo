package middleware

import (
	"astigo/internal/domain/handler"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	handler handler.IAuthHandler
}

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

func NewAuthMiddleware(authHandler handler.IAuthHandler) *AuthMiddleware {
	return &AuthMiddleware{
		handler: authHandler,
	}
}
