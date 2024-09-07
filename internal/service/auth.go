package service

import (
	"context"
	"github.com/passionde/user-segmentation-service/internal/repo"
	"github.com/passionde/user-segmentation-service/pkg/secure"
)

type AuthService struct {
	authRepo repo.Auth
	secure   secure.APISecure
}

func NewAuthService(authRepo repo.Auth, secure secure.APISecure) *AuthService {
	return &AuthService{
		authRepo: authRepo,
		secure:   secure,
	}
}

func (a *AuthService) TokenExist(ctx context.Context, token string) (int, error) {
	return a.authRepo.TokenExist(ctx, a.secure.Hash(token))
}

func (a *AuthService) GenerateToken(ctx context.Context) (int, string, error) {
	token := a.secure.GenerateKey()
	id, err := a.authRepo.WriteToken(ctx, a.secure.Hash(token))
	return id, token, err
}
