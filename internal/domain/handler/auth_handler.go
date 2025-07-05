package handler

import (
	"astigo/internal/domain/model"
	"context"
	"github.com/coreos/go-oidc"
)

type IAuthHandler interface {
	VerifyToken(ctx context.Context, token string) (*oidc.IDToken, error)
	GetClaims(idToken *oidc.IDToken) (*model.Claims, error)
}
