package middleware

import (
	"net/http"
	"strings"

	"github.com/TancelinMazzotti/astigo/internal/domain/model"
	"github.com/TancelinMazzotti/astigo/internal/domain/service"

	"github.com/gin-gonic/gin"
)

var _ IAuthMiddleware = (*AuthMiddleware)(nil)

type IAuthMiddleware interface {
	Middleware(ctx *gin.Context)
	CheckRealmMiddleware(roles []string) func(ctx *gin.Context)
	CheckResourceRoleMiddleware(resource string, roles []string) func(ctx *gin.Context)
}

// AuthMiddleware handles authentication by validating JWT tokens in incoming requests.
// It uses a provided implementation of IAuthService for token verification and claims extraction.
type AuthMiddleware struct {
	authService service.IAuthService
}

// Middleware is a Gin middleware function that validates JWT Authorization headers for protected routes.
func (m *AuthMiddleware) Middleware(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header"})
		return
	}

	token := parts[1]

	idToken, err := m.authService.VerifyToken(ctx, token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, err := m.authService.GetClaims(idToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid payload"})
		return
	}

	ctx.Set("claims", claims)

	ctx.Next()

}

// CheckRealmMiddleware checks if a user's JWT claims include at least one of the specified realm roles and authorizes accordingly.
func (m *AuthMiddleware) CheckRealmMiddleware(roles []string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		claimsCtx, exists := ctx.Get("claims")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims type"})
			return
		}

		claims, ok := claimsCtx.(*model.Claims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims type"})
			return
		}

		for _, role := range roles {
			if claims.HasRealmRole(role) {
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "forbidden"})
	}
}

// CheckResourceRoleMiddleware validates whether the user's claims include the required roles for a specific resource.
func (m *AuthMiddleware) CheckResourceRoleMiddleware(resource string, roles []string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		claimsCtx, exists := ctx.Get("claims")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims type"})
			return
		}

		claims, ok := claimsCtx.(*model.Claims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims type"})
			return
		}

		for _, role := range roles {
			if claims.HasResourceRole(resource, role) {
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "forbidden"})
	}
}

// NewAuthMiddleware creates and returns an instance of AuthMiddleware using the provided IAuthService for authentication.
func NewAuthMiddleware(authService service.IAuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}
