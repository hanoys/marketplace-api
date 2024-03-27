package service

import (
	"context"
	"github.com/hanoys/marketplace-api/internal/domain"
)

type UsersRepository interface {
	Create(ctx context.Context, login string, password string) (domain.User, error)
	FindByLogin(ctx context.Context, login string) (domain.User, error)
}

type AdvertisementsRepository interface {
	Create(ctx context.Context, userID int, title string, body string, imageURL string, price float64) (domain.Advertisement, error)
	GetAdvertisements(ctx context.Context, userID int, pageNum int, sort domain.SortType, dir domain.DirectionType) ([]domain.AdvertisementEntry, error)
	FindAll(ctx context.Context) ([]domain.Advertisement, error)
}

type Repositories struct {
	UsersRepository
	AdvertisementsRepository
}
