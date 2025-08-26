package service

import (
	"context"
	"fmt"

	"github.com/TancelinMazzotti/astigo/internal/domain/model"

	"github.com/coreos/go-oidc"
	"go.uber.org/zap"
)

var _ IAuthService = (*AuthService)(nil)

type IAuthService interface {
	VerifyToken(ctx context.Context, token string) (*oidc.IDToken, error)
	GetClaims(idToken *oidc.IDToken) (*model.Claims, error)
}

// AuthService handles authentication tasks using OpenID Connect provider.
// It verifies ID tokens and extracts claims.
// Dependencies include a logger, an OIDC provider, and a token verifier.
type AuthService struct {
	logger   *zap.Logger
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
}

// VerifyToken verifies the provided ID token using the configured OIDC verifier and returns the parsed token or an error.
func (s *AuthService) VerifyToken(ctx context.Context, token string) (*oidc.IDToken, error) {
	idToken, err := s.verifier.Verify(ctx, token)
	if err != nil {
		s.logger.Debug("failed to verify token", zap.Error(err))
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}
	return idToken, nil
}

// GetClaims extracts claims from the provided ID token and returns them as a Claims object or an error if it fails.
func (s *AuthService) GetClaims(idToken *oidc.IDToken) (*model.Claims, error) {
	claims := &model.Claims{}
	if err := idToken.Claims(&claims); err != nil {
		s.logger.Debug("failed to get claims", zap.Error(err))
		return nil, fmt.Errorf("failed to get claims: %w", err)
	}
	return claims, nil
}

// NewAuthService initializes and returns a new AuthService instance with the provided logger, OIDC provider, and client ID.
func NewAuthService(logger *zap.Logger, provider *oidc.Provider, clientId string) *AuthService {
	return &AuthService{
		logger:   logger,
		provider: provider,
		verifier: provider.Verifier(&oidc.Config{ClientID: clientId}),
	}
}
