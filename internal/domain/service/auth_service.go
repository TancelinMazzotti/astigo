package service

import (
	"astigo/internal/domain/handler"
	"astigo/internal/domain/model"
	"context"
	"fmt"
	"github.com/coreos/go-oidc"
)

var _ handler.IAuthHandler = (*AuthService)(nil)

type AuthConfig struct {
	ClientID string `mapstructure:"client_id"`
	Issuer   string `mapstructure:"issuer"`
}

type AuthService struct {
	config   AuthConfig
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
}

func (s *AuthService) VerifyToken(ctx context.Context, token string) (*oidc.IDToken, error) {
	idToken, err := s.verifier.Verify(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}
	return idToken, nil
}

func (s *AuthService) GetClaims(idToken *oidc.IDToken) (*model.Claims, error) {
	claims := &model.Claims{}
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to get claims: %w", err)
	}
	return claims, nil
}

func NewAuthService(config AuthConfig) *AuthService {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, config.Issuer)
	if err != nil {
		panic("Failed to initialize Keycloak provider: " + err.Error())
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: config.ClientID})

	return &AuthService{
		config:   config,
		provider: provider,
		verifier: verifier,
	}
}
