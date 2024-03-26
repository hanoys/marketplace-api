package service

import (
	"context"
	"github.com/hanoys/marketplace-api/internal/domain"
)

type Users interface {
	Create(ctx context.Context, login string, password string) (domain.User, error)
}

type Advertisements interface {
	Create(ctx context.Context, userID int, title string, body string, imageURL string, price float64) (domain.Advertisement, error)
	FindAll(ctx context.Context) ([]domain.Advertisement, error)
}

type Services struct {
	Users
	Advertisements
}

func NewServices(repositories *domain.Repositories) *Services {
	return &Services{Users: NewUserService(repositories),
		Advertisements: NewAdvertisementService(repositories)}
}
