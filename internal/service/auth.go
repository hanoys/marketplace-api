package service

import (
	"context"
	"github.com/hanoys/marketplace-api/auth"
)

type AuthorizationService struct {
	repositories  *Repositories
	tokenProvider *auth.Provider
}

func NewAuthorizationService(repositories *Repositories, tokenProvider *auth.Provider) *AuthorizationService {
	return &AuthorizationService{repositories: repositories,
		tokenProvider: tokenProvider}
}

// TODO: change errors
func (a *AuthorizationService) LogIn(ctx context.Context, login string, password string) (*auth.TokenPair, error) {
	user, err := a.repositories.UsersRepository.FindByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	tokenPayload, err := a.tokenProvider.NewPayload(user.ID)
	if err != nil {
		return nil, err
	}

	session, err := a.tokenProvider.NewSession(ctx, tokenPayload)
	if err != nil {
		return nil, err
	}

	return session.Tokens, nil
}

func (a *AuthorizationService) LogOut(ctx context.Context, tokenString string) error {
	err := a.tokenProvider.CloseSession(ctx, tokenString)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthorizationService) RefreshToken(ctx context.Context, refreshTokenString string) (*auth.TokenPair, error) {
	session, err := a.tokenProvider.RefreshSession(ctx, refreshTokenString)
	if err != nil {
		return nil, err
	}

	return session.Tokens, nil
}

func (a *AuthorizationService) VerifyToken(ctx context.Context, tokenString string) (int, error) {
	payload, err := a.tokenProvider.VerifyToken(ctx, tokenString)
	if err != nil {
		return 0, err
	}

	return payload.UserID, err
}
