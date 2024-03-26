package service

import (
	"context"
	"github.com/hanoys/marketplace-api/internal/domain"
)

type UserService struct {
	repositories *domain.Repositories
}

func NewUserService(repositories *domain.Repositories) *UserService {
	return &UserService{repositories: repositories}
}

func (u *UserService) Create(ctx context.Context, login string, password string) (domain.User, error) {
	return u.repositories.Users.Create(ctx, login, password)
}
