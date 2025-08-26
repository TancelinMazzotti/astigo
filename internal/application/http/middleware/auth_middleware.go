package middleware

import (
	"astigo/internal/domain/model"
	"astigo/internal/domain/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid payload"})
		return
	}

	c.Set("claims", claims)

	c.Next()

}

// CheckRealmMiddleware checks if a user's JWT claims include at least one of the specified realm roles and authorizes accordingly.
func (m *AuthMiddleware) CheckRealmMiddleware(roles []string) func(c *gin.Context) {
	return func(c *gin.Context) {
		claimsCtx, exists := c.Get("claims")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims type"})
			return
		}

		claims, ok := claimsCtx.(*model.Claims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims type"})
			return
		}

		for _, role := range roles {
			if claims.HasRealmRole(role) {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "forbidden"})
	}
}

// CheckResourceRoleMiddleware validates whether the user's claims include the required roles for a specific resource.
func (m *AuthMiddleware) CheckResourceRoleMiddleware(resource string, roles []string) func(c *gin.Context) {
	return func(c *gin.Context) {
		claimsCtx, exists := c.Get("claims")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims type"})
			return
		}

		claims, ok := claimsCtx.(*model.Claims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims type"})
			return
		}

		for _, role := range roles {
			if claims.HasResourceRole(resource, role) {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "forbidden"})
	}
}

// NewAuthMiddleware creates and returns an instance of AuthMiddleware using the provided IAuthService for authentication.
func NewAuthMiddleware(authHandler service.IAuthService) *AuthMiddleware {
	return &AuthMiddleware{
		handler: authHandler,
	}
}
