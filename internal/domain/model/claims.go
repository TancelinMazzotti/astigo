package model

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"slices"
)

const (
	ClaimsContextKey = "claims"
)

type Claims struct {
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

func (c *Claims) HasRealmRole(role string) bool {
	return slices.Contains(c.RealmAccess.Roles, role)
}

func (c *Claims) HasResourceRole(resource string, role string) bool {
	if c.ResourceAccess == nil {
		return false
	}

	if _, ok := c.ResourceAccess[resource]; !ok {
		return false
	}

	return slices.Contains(c.ResourceAccess[resource].Roles, role)
}

func SetClaimsInContext(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, ClaimsContextKey, claims)
}

func GetClaimsInContext(ctx context.Context) (*Claims, error) {
	value := ctx.Value(ClaimsContextKey)
	if value == nil {
		return nil, fmt.Errorf("no claims found in context")
	}

	claims, ok := value.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid claims type in context")
	}

	return claims, nil
}
