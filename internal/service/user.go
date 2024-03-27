package service

import (
	"context"
	"fmt"
	"github.com/hanoys/marketplace-api/internal/domain"
)

type UserService struct {
	repositories *Repositories
}

func NewUserService(repositories *Repositories) *UserService {
	return &UserService{repositories: repositories}
}

func (a *UserService) SignUp(ctx context.Context, login string, password string) (domain.User, error) {
	_, err := a.repositories.UsersRepository.FindByLogin(ctx, login)
	// TODO: to distinguish an error about not found user and db error
	if err == nil {
		return domain.User{}, fmt.Errorf("user already exists")
	}

	user, err := a.repositories.UsersRepository.Create(ctx, login, password)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
