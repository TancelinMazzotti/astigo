package middleware

import (
	"context"
	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Email             string `json:"email"`
	PreferredUsername string `json:"preferred_username"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	RealmAccess       struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	ResourceAccess map[string]struct {
		Roles []string `json:"roles"`
	} `json:"resource_access"`
}

type AuthMiddleware struct {
	ctx      context.Context
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
	clientID string
	issuer   string
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

	idToken, err := m.verifier.Verify(m.ctx, token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var claims CustomClaims
	if err := idToken.Claims(&claims); err != nil {

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse claims"})
		return
	}

	c.Set("claims", claims)

	c.Next()

}

func NewAuthMiddleware(issuer string, clientID string) *AuthMiddleware {

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		panic("Failed to initialize Keycloak provider: " + err.Error())
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: clientID})

	return &AuthMiddleware{
		ctx:      ctx,
		provider: provider,
		verifier: verifier,
		clientID: clientID,
		issuer:   issuer,
	}
}
